# Implementation Plan — Admin Authentication & Real Data Integration

We will implement a secure, professional administrative authentication flow and connect real database statistics to the Admin Dashboard overview panels. 

---

## User Review Required

> [!IMPORTANT]
> **Key Architecture & Security Shifts**
> 1. **Default Application Entry**: Normal citizens and admins now land exclusively on the citizen portal (`/`) upon starting the project. There is **no automatic redirection** to the admin panel on startup.
> 2. **Dedicated Routes**:
>    * `/admin-dashboard`: A clean, isolated **Admin Login Page** with dynamic validation.
>    * `/admin/dashboard`: The actual **Protected Admin Dashboard**.
> 3. **Route Protection**: Anyone attempting to access `/admin/dashboard` (or any sub-route under `/admin/`) without an active administrator session is immediately redirected back to `/admin-dashboard` (instead of opening citizen modals).
> 4. **Dynamic Database Overview**: We replace mock stats with real live counts of total schemes, users, applications (pending, approved, rejected), dynamic applicant listings, and activity feeds.

---

## Proposed Changes

We will group our work into the following components:

### 1. Go Backend (API Analytics Integration)

#### [MODIFY] [admin_handlers.go](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/backend/internal/handlers/admin_handlers.go)
* Modify the `GetAdminAnalyticsHandler` API to query dynamic operational database metrics:
  * Count total applications submitted (`total_applications`).
  * Count applications grouped by status (`pending_applications`, `approved_applications`, `rejected_applications`).
  * Query the last 5 recent applications (applicant name, scheme title, status, applied date).
  * Format these metrics and pack them into the JSON output payload under `total_applications`, `pending_applications`, `approved_applications`, `rejected_applications`, and `recent_applications`.

---

### 2. Vue Frontend (Routing & Authentication Isolation)

#### [MODIFY] [router/index.js](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/frontend/src/router/index.js)
* Register a new public route `/admin-dashboard` pointing to the new `AdminLoginView.vue`.
* Rename the path of the `/admin` view to `/admin/dashboard`.
* Update the `beforeEach` navigation guard:
  * If a user tries to access a path starting with `/admin/` (such as `/admin/dashboard`), verify `authStore.isAdmin`.
  * If they are not logged in as an admin, block the transition and redirect them directly to the dedicated `/admin-dashboard` login portal.

#### [MODIFY] [App.vue](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/frontend/src/App.vue)
* Delete the automatic startup/reload `watch(() => authStore.isAdmin)` that directly pushed `/admin`.
* Update `handleTabChange` to map `admin` to the new path `/admin/dashboard`.
* Update the conditional layout flags (`v-if="$route.name !== 'admin'"` ) to also hide citizen `<Header>` and `<Hero>` components when the route name is `'admin-dashboard'` (so both the login page and dashboard are fully isolated).

#### [NEW] [AdminLoginView.vue](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/frontend/src/views/AdminLoginView.vue)
* Build a beautiful, professional, and responsive admin login view strictly following flat-design principles:
  * curating premium dual typography (`Outfit` for headings, `Inter` for inputs).
  * 0.5px borders (`0.5px solid rgba(0,0,0,0.08)`) with custom color variables.
  * Fields for Administrative Email and Password with proper error message alerts if validation fails (e.g. "Invalid admin credentials").
  * Secure auth action: logs in with `authStore.loginUser()`. On success, checks `authStore.isAdmin`. If true, redirects to `/admin/dashboard`. If false (a citizen logged in on the admin page), logs them out immediately and throws an error "Access Denied: Administrative privileges required."

#### [MODIFY] [AdminSidebar.vue](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/frontend/src/components/admin/AdminSidebar.vue)
* Update the bottom logout button trigger:
  * Clear token session using `authStore.logoutUser()`.
  * Redirect the admin directly to the Admin Login Page `/admin-dashboard`!

---

### 3. Vue Frontend (Real Dashboard Overview & Stats)

#### [MODIFY] [OverviewTab.vue](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/frontend/src/components/admin/OverviewTab.vue)
* Bind the 4 overview statistics cards dynamically:
  * Total Schemes (`sa.total_schemes`)
  * Registered Users (`sa.total_users`)
  * Total Applications (`sa.total_applications || 0`) — instead of Categories.
  * Pending Approvals (`sa.pending_applications || 0`) — instead of Expiring Soon.
* Bind the left table to **"Recent Applications"** instead of hardcoded schemes, showing real database applications mapping `props.analytics.recent_applications` with columns `Applicant`, `Scheme`, and `Status` badges.
* Bind `recentActivity` dynamically using `props.analytics.recent_activity` logs from the database.

#### [MODIFY] [AnalyticsTab.vue](file:///d:/Team-Ai/Project-001/yojana-portal/majhiGov/frontend/src/components/admin/AnalyticsTab.vue)
* Bind the 4 analytics cards dynamically to show real database-driven results:
  * Total Applications (`props.analytics.total_applications || 0`)
  * Pending Approvals (`props.analytics.pending_applications || 0`)
  * Approved Applications (`props.analytics.approved_applications || 0`)
  * Rejected Applications (`props.analytics.rejected_applications || 0`)

---

## Verification Plan

### Automated Build Verification
1. Run backend tests to ensure the server starts cleanly:
   ```powershell
   go test ./...
   ```
2. Compile the Vue 3 frontend in production mode to confirm 100% build validity:
   ```powershell
   npm run build
   ```

### Manual Verification
1. **Startup Check**: Run the application and verify it opens the citizen home page at `/` (and does not redirect).
2. **Access Protection**: Try typing `http://localhost:5173/#/admin/dashboard` directly in an incognito window. Verify it redirects you to the Admin Login Page at `/admin-dashboard` instead of showing any content.
3. **Invalid Login**: Go to `/admin-dashboard`, enter incorrect details (e.g. `admin@gov.in` / `wrongpassword`), and verify the message "Invalid credentials" is shown.
4. **Valid Login**: Enter `admin@gov.in` / `admin123` on `/admin-dashboard`. Verify you are immediately logged in, a secure session is stored, and you are redirected to the Admin Dashboard.
5. **Real-time Overview Data**: Confirm that all numbers on the stats cards, recent applications table, and activity feed correspond to actual database counts rather than static placeholders.
6. **Logout Flow**: Click the exit/logout button inside the Admin Sidebar. Verify it clears your local token and redirects you back to `/admin-dashboard`.
