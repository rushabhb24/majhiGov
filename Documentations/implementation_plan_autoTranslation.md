# Yojana Portal: Admin Enhancements & Automatic Translation Plan

This implementation plan details the architectural and UI upgrades for the majhiGov Yojana Portal's Admin Dashboard. We will add a dedicated administrative profile manager, navigation and portal return paths, and a translation proxy allowing admins to create bilingual/trilingual schemes automatically with manual correction support.

---

## Technical Overview & Proposed Architecture

### 1. Database-Backed Administrative Profiles
We will extend the `user_profiles` schema to store profile photos (`avatar_url`). The profile editing section will let the logged-in administrator update:
- Personal credentials in the `users` table: **Email**, **Phone**, and **Password** (secured via `bcrypt` hashing with full uniqueness checks).
- Demographic details in the `user_profiles` table: **Full Name** and **Avatar URL**.
- We will execute both updates atomically using a SQL database transaction (`tx`) to prevent partial data synchronization.

### 2. Side-by-Side Automatic Translation Proxy
To translate newly inputted schemes from English into Hindi and Marathi effortlessly, we will build a server-side proxy route:
`GET /api/translate?q={text}&target={hi|mr}`
This proxy will query Google's public translation API safely, bypassing frontend CORS blocks, parsing multi-sentence inputs, and returning translated strings. On the frontend, a blur/change event on English title/description/FAQ inputs will automatically invoke this proxy, populating Hindi and Marathi inputs dynamically while keeping them fully interactive and editable.

### 3. Glassmorphic Profile Interface & Dashboard Routing
- **Sidebar Integration**: We will introduce a premium **"Return to Portal"** nav button with an interactive Tabler icon (`ti-arrow-back-up`) under a new sidebar section label. We will also add a **"Profile"** nav link.
- **`ProfileTab.vue` [NEW]**: A stunning settings dashboard matching our aesthetic system, featuring:
  - An interactive **Avatar Selection Deck** loaded with modern administrative illustration cards.
  - A secure credentials form (Name, Email, Phone, Password).
  - Toast confirmation notifications indicating updates.

---

## User Review Required

> [!IMPORTANT]
> **Database Column Migration**:
> Adding the `avatar_url` column to `user_profiles` will be fully automated in the `db.go` startup migrations using:
> `ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500) DEFAULT '';`
> This maintains complete backwards compatibility for existing user databases.

> [!TIP]
> **Translation Performance Optimization**:
> To ensure typing is responsive and light, translations will trigger **only** when the user finishes writing in an English input field and clicks out (focus loss via the `@blur` event). This prevents redundant HTTP requests during active typing.

---

## Open Questions

> [!WARNING]
> **1. Avatar Storage Preference**:
> We propose rendering a gorgeous card deck of **12 pre-styled SVG illustrations** (avatars representing officers, administrators, and avatars of varying styles) which the admin can select with one-click. This avoids complex file upload issues and provides a beautiful, modern look. Does this pre-seeded deck approach suit your vision, or do you prefer a text input field where you paste any arbitrary image URL?
> *(Our recommendation is the pre-seeded SVG/illustration deck for maximum premium feel).*

---

## Proposed Changes

### 1. Database Schema & Startup Migrations
#### [MODIFY] [db.go](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/backend/internal/db/db.go)
- Add a table schema alteration to `runMigrations()`:
  ```sql
  ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500) DEFAULT '';
  ```

#### [MODIFY] [models.go](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/backend/internal/models/models.go)
- Extend `UserProfile` struct to map `avatar_url`, `email` and `phone` joined fields:
  ```go
  AvatarURL string `json:"avatar_url" db:"avatar_url"`
  Email     string `json:"email" db:"email"` // Joined field
  Phone     string `json:"phone" db:"phone"` // Joined field
  ```

---

### 2. Backend Routing & Handlers
#### [MODIFY] [main.go](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/backend/cmd/server/main.go)
- Register the secure translation proxy endpoint:
  ```go
  mux.Handle("/api/translate", middleware.AuthMiddleware(http.HandlerFunc(handlers.TranslateHandler)))
  ```

#### [MODIFY] [handlers.go](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/backend/internal/handlers/handlers.go)
- **`TranslateHandler` [NEW]**: Proxies calls to Google Translate's single translation API, joining multi-sentence arrays cleanly.
- **`GetUserProfileHandler`**: JOIN `users` to select and scan `email`, `phone`, and `avatar_url`.
- **`UpdateUserProfileHandler`**: Begin a database transaction (`tx`). Perform uniqueness checks for modified email/phone. Hashing the password with bcrypt if provided. Perform update queries on `users` and `user_profiles`. Commit safely and return updated fields.
- **`LoginHandler`**: Join the `avatar_url` into the queried profile.

---

### 3. Frontend Authentication Stores
#### [MODIFY] [auth.js](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/frontend/src/stores/auth.js)
- Extend `updateProfile` action to accept comprehensive admin details (email, phone, password, avatar_url) and securely pass them to the PUT `/api/user/profile` endpoint.

---

### 4. Admin Dashboard Interfaces
#### [MODIFY] [AdminSidebar.vue](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/frontend/src/components/admin/AdminSidebar.vue)
- Add the **"PORTAL"** sidebar section label with the **"Return to Portal"** navigation item (`ti ti-arrow-back-up` icon).
- Add the **"Admin Profile"** item (`ti ti-user` icon) under the **"SYSTEM"** label.
- Bind the avatar initials circle inside the sidebar footer to dynamically render the administrator's custom selected avatar image if available.

#### [MODIFY] [AdminView.vue](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/frontend/src/views/AdminView.vue)
- Add `ProfileTab.vue` component import and render it inside the page content matching `activeTab === 'profile'`.
- Implement `translateField(text, fieldKey)`, `translateDocRow(doc)`, `translateFaqRow(faq)` calling our backend proxy.
- Bind blur `@blur` event handlers to English Scheme identity fields (Title, Description), FAQ Questions/Answers, and Document names in the **Add New Scheme** form overlay.

#### [NEW] [ProfileTab.vue](file:///t:/MKCL-Office/Projects/Projects/yojana-portal/frontend/src/components/admin/ProfileTab.vue)
- Create a beautiful profile edit tab:
  - Interactive **Avatar Deck Card Selector** (featuring 8 curated vector avatars).
  - Form inputs: Full Name, email address, telephone, and a secure password modification field.
  - Submit trigger calling `authStore.updateProfile()`.

---

## Verification Plan

### Automated Build Validation
- Run Vite local development validation:
  `npm run build`
- Run Go backend builds to assert server compilation:
  `go build .\cmd\server\`

### Manual Functional Validation
1. **Sidebar Navigation**:
   - Verify clicking the **"Return to Portal"** button correctly redirects to the main route (`/`).
   - Verify clicking **"Admin Profile"** successfully mounts the profile editing view.
2. **Admin Profile Editing**:
   - Edit the Name, Email, and Phone fields. Save and verify that details sync in the DB and display in the sidebar footer.
   - Choose a different avatar illustration card. Ensure the avatar preview updates instantly in both the form and the sidebar.
3. **Translation Engine Verification**:
   - Open **"Add New Scheme"**.
   - Input `"PM Kisan Welfare Scheme"` into the English Title. Press Tab (triggering focus blur). Verify that Hindi and Marathi title inputs populate automatically with their translations.
   - Change a translated character to check that fields remain fully editable for manual adjustments.
