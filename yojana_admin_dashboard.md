# Yojana Portal — Admin Dashboard Complete Prompt

---

## ROLE & CONTEXT

You are an expert full-stack developer building a production-ready
Admin Dashboard for "Yojana Portal" — a Government Schemes Portal
for citizens of India. Build the complete admin panel as a single
self-contained HTML file with embedded CSS and JavaScript.
No external frameworks. No build tools. Just one file that works
when opened in a browser.

---

## TECH & STACK REQUIREMENTS

- Single HTML file (index.html)
- Vanilla CSS with CSS variables for theming
- Vanilla JavaScript (no jQuery, no React, no Vue)
- Tabler Icons via CDN for all icons:
  https://cdn.jsdelivr.net/npm/@tabler/icons-webfont@latest/tabler-icons.min.css
- Google Fonts via CDN:
  https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600&display=swap
- No other external dependencies

---

## LAYOUT STRUCTURE

The dashboard uses a fixed two-panel layout:

```
┌─────────────────────────────────────────────────┐
│  SIDEBAR (220px fixed left)  │  MAIN AREA        │
│                              │  ┌─────────────┐  │
│  Logo + App Name             │  │  TOP BAR    │  │
│  ─────────────────           │  └─────────────┘  │
│  Navigation Menu             │  │  PAGE       │  │
│  ─────────────────           │  │  CONTENT    │  │
│  Admin User Info             │  │  (scrolls)  │  │
└──────────────────────────────┴──────────────────-┘
```

- Sidebar: 220px wide, fixed height 100vh, does NOT scroll
- Main area: fills remaining width, has fixed topbar + scrollable content
- Total height: 100vh (full screen)
- Font: Plus Jakarta Sans throughout

---

## COLOR SYSTEM (CSS Variables)

Define these exact CSS variables in :root:

```css
:root {
  --primary: #1a3a6b;          /* Deep navy blue — sidebar bg, primary buttons */
  --primary-light: #e8eef8;    /* Light blue — badge bg, hover states */
  --accent: #f97316;           /* Orange — CTA buttons, badges, highlights */
  --accent-light: #fff4ed;     /* Light orange — accent badge bg */
  --success: #16a34a;          /* Green — active status, success states */
  --success-bg: #f0fdf4;       /* Light green — success badge bg */
  --danger: #dc2626;           /* Red — delete, danger, inactive */
  --danger-bg: #fef2f2;        /* Light red — danger badge bg */
  --warning: #d97706;          /* Amber — expiring, warning states */
  --warning-bg: #fffbeb;       /* Light amber — warning badge bg */
  --border: rgba(0,0,0,0.08); /* All borders */
  --text: #0f172a;             /* Primary text */
  --text2: #64748b;            /* Secondary/muted text */
  --bg: #ffffff;               /* Card/surface bg */
  --bg2: #f8fafc;              /* Subtle surface bg */
  --bg3: #f1f5f9;              /* Page/outer bg */
  --radius: 8px;               /* Standard border radius */
  --radius-lg: 12px;           /* Large border radius for cards */
}
```

---

## SIDEBAR — EXACT SPECIFICATION

### Brand Area (top of sidebar)
- Background: var(--primary)
- Logo box: 32x32px, background var(--accent), border-radius 8px
  Icon inside: ti-building-bank (Tabler icon), white, 18px
- App name: "Yojana Portal" in white, 13px, weight 500
- Subtitle below it: "Admin Panel" in rgba(255,255,255,0.5), 11px

### Navigation Menu
Sections with uppercase labels in rgba(255,255,255,0.4), 10px, letter-spacing 0.08em

SECTION: "MAIN"
  Item 1: ti-layout-dashboard  →  "Overview"        [default active]
  Item 2: ti-files             →  "Schemes"         [badge: "142" in orange]
  Item 3: ti-grid-dots         →  "Categories"

SECTION: "MANAGEMENT"
  Item 4: ti-users             →  "Users"
  Item 5: ti-checklist         →  "Eligibility Rules"
  Item 6: ti-bell              →  "Notifications"   [badge: "3" in orange]

SECTION: "SYSTEM"
  Item 7: ti-chart-bar         →  "Analytics"
  Item 8: ti-settings          →  "Settings"

### Nav Item Styling
- Normal: rgba(255,255,255,0.7) text, transparent bg
- Hover: rgba(255,255,255,0.08) bg, white text
- Active: rgba(255,255,255,0.15) bg, white text, icon color = var(--accent)
- Height: ~36px, border-radius 6px, padding 9px 10px
- Icon: 16px, 18px wide fixed (for alignment)

### Badge Styling (on nav items)
- Background: var(--accent)
- Color: white
- Font-size: 10px
- Padding: 2px 6px
- Border-radius: 10px
- Float to right (margin-left: auto)

### Footer (bottom of sidebar)
- Thin top border: 0.5px solid rgba(255,255,255,0.1)
- Shows admin avatar circle (30x30, initials "SA")
  bg: rgba(255,255,255,0.15), text white, 12px
- Name: "Super Admin", rgba(255,255,255,0.8), 12px
- Role: "Administrator", rgba(255,255,255,0.4), 11px

---

## TOP BAR — EXACT SPECIFICATION

- Height: 56px
- Background: white
- Border-bottom: 0.5px solid var(--border)
- Padding: 0 20px
- Flex row with: Page Title | Search Bar | Bell Button | Primary Action Button

### Page Title
- Font-size: 15px, weight 500
- Changes dynamically when sidebar nav item is clicked
- flex: 1 (takes remaining space)

### Search Bar
- Width: 200px
- Background: var(--bg2)
- Border: 0.5px solid var(--border), border-radius 6px
- Padding: 6px 10px
- Contains: ti-search icon (15px, muted) + text input (no border, transparent bg)
- Placeholder: "Search schemes, users..."

### Bell Button
- Icon: ti-bell
- Small orange dot positioned top-right of button (notification indicator)
- Dot: 7px circle, background var(--accent)

### Primary Action Button
- Background: var(--primary), color white
- Padding: 7px 14px, border-radius 6px, font-size 13px
- Contains: ti-plus icon + dynamic label text
- Label changes per page:
  Overview       → "Add Scheme"
  Schemes        → "Add New Scheme"
  Categories     → "Add Category"
  Users          → "Add Admin"
  Eligibility    → "Save Rules"
  Notifications  → "Send Notification"
  Analytics      → "Export Report"

---

## PAGES — DETAILED SPECIFICATION

### HOW PAGES WORK
- All 7 page content divs exist in the DOM simultaneously
- Only one has display:block at a time (others display:none)
- Clicking a sidebar nav item:
  1. Hides all pages
  2. Shows the clicked page
  3. Removes .active from all nav items, adds to clicked one
  4. Updates page title text
  5. Updates primary action button text
- No page reload. Pure JavaScript DOM switching.

---

### PAGE 1: OVERVIEW (default visible)

#### Stats Row — 4 cards in a grid (equal width columns)
Each stat card: white bg, 0.5px border, border-radius 8px, padding 14px 16px

Card 1 — Total Schemes
  Icon box: 36x36, bg var(--primary-light), color var(--primary), icon ti-files, border-radius 8px
  Number: "142", 22px, weight 500
  Label: "Total Schemes", 12px muted
  Change: "+8 this month" in green with ti-trending-up icon

Card 2 — Registered Users
  Icon box: bg var(--success-bg), color var(--success), icon ti-users
  Number: "18,430"
  Label: "Registered Users"
  Change: "+1,204 this week" in green

Card 3 — Categories
  Icon box: bg var(--accent-light), color var(--accent), icon ti-grid-dots
  Number: "10"
  Label: "Categories"
  Change: "All active" in green with ti-check icon

Card 4 — Expiring Soon
  Icon box: bg var(--danger-bg), color var(--danger), icon ti-clock
  Number: "7"
  Label: "Expiring Soon"
  Change: "Within 30 days" in red with ti-alert-circle icon

#### Warning Alert Banner
Below stats, full width:
- Background: var(--warning-bg)
- Border: 0.5px solid #fbbf24
- Border-radius: 8px
- Padding: 10px 14px
- Icon: ti-alert-triangle in var(--warning), 16px
- Text: "7 schemes are expiring within 30 days. Review and update their deadlines."
- Right side: "View all" text in underline, cursor pointer

#### Two-Column Section (equal width)

LEFT CARD — "Recent Schemes Added"
Card header: title left, "View all →" link right (ti-arrow-right, primary color)
Inside: mini table with columns: Scheme | Type | Status
No outer border on table. Row hover: bg2.
Rows (4 total):
  Row 1: "PM Kisan Samman" / subtitle "Agriculture" | badge "Central" (primary-light bg) | badge "Active" (green)
  Row 2: "NSP Scholarship" / subtitle "Education" | "Central" | "Active"
  Row 3: "Ladli Behna Yojana" / subtitle "Women" | "State" (accent-light bg, orange) | "Expiring" (warning-bg, amber)
  Row 4: "Mudra Loan" / subtitle "Business" | "Central" | "Active"

RIGHT CARD — "Schemes by Category"
Card header: "Schemes by Category"
Inside: list of category rows, each row shows:
  [emoji icon] [Category Name] → [progress bar 80px wide] [count]
  Progress bar: thin (4px), border-radius 2px, bg var(--bg2)
  Fill: var(--primary)

Rows:
  🌾 Farmers      → bar 85% full  → "34"
  🎓 Students     → bar 70%       → "28"
  👩 Women        → bar 55%       → "22"
  💼 Business     → bar 45%       → "18"
  👴 Senior Citizens → bar 35%    → "14"
  ♿ Disabled     → bar 22%       → "9"

#### Activity Feed Card (full width)
Title: "Recent Activity"
4 activity items, each has:
- Colored dot (8px circle, left side)
- Activity text (13px)
- Time ago (11px, muted)
Separated by thin border-bottom lines

Items:
  🟢 "New scheme PM Vishwakarma Yojana added to Business category" — "2 hours ago"
  🟠 "Scheme Ladli Behna deadline updated to 31 March 2025" — "5 hours ago"
  🔵 "User Ramesh Kumar registered from Rajasthan" — "1 day ago"
  🟢 "Category Senior Citizens updated with 2 new schemes" — "2 days ago"

---

### PAGE 2: SCHEMES

#### Filter Row (above table)
Right-aligned, 3 select dropdowns:
  Dropdown 1: "All Categories" / Farmers / Students / Women
  Dropdown 2: "All Types" / Central / State
  Dropdown 3: "All Status" / Active / Inactive / Expiring
Each: border 0.5px, border-radius 6px, font-size 13px, padding 7px 10px

#### Schemes Table (full width card)
Columns: Scheme Name | Category | Type | Deadline | Status | Actions

Rows (5 total):
  1. PM Kisan Samman Nidhi
     subtitle: "₹6,000/year to farmers"
     Category: 🌾 Farmers | Type: Central | Deadline: Ongoing | Status: Active

  2. NSP Post-Matric Scholarship
     subtitle: "Scholarship for SC/ST students"
     Category: 🎓 Students | Type: Central | Deadline: 31 Mar 2025 | Status: Active

  3. Ladli Behna Yojana
     subtitle: "₹1,250/month to women"
     Category: 👩 Women | Type: State - MP | Deadline: 28 Feb 2025 | Status: Expiring

  4. PM Mudra Yojana
     subtitle: "Loans up to ₹10 lakh"
     Category: 💼 Business | Type: Central | Deadline: Ongoing | Status: Active

  5. Indira Gandhi Pension Yojana
     subtitle: "Monthly pension for elderly"
     Category: 👴 Senior Citizens | Type: Central | Deadline: Ongoing | Status: Inactive

#### Action Buttons (per row, 3 buttons)
Each button: 28x28px, border-radius 6px, border 0.5px, icon 15px
  Button 1: ti-eye (view) — hover: bg2
  Button 2: ti-edit (edit) — hover: bg2
  Button 3: ti-trash (delete) — hover: danger-bg, danger color border

#### Table Styling
- Header: bg2, uppercase, 11px, muted, letter-spacing 0.05em
- Row hover: bg2 background
- Border-bottom on each row (not last)
- Scheme name: 500 weight. Subtitle: 11px muted below

---

### PAGE 3: CATEGORIES

Two equal columns side by side:

LEFT CARD — "Add New Category" (form)
Fields:
  1. Label "Category Name (English)" → text input, placeholder "e.g. Farmers"
  2. Two-col row:
     - Label "Hindi Name" → input, placeholder "किसान"
     - Label "Marathi Name" → input, placeholder "शेतकरी"
  3. Label "Icon (Emoji)" → input, placeholder "🌾"
  4. Label "Description" → input, placeholder "Schemes for farmers and agriculture"
  5. Submit button: "Add Category" with ti-plus icon
     bg var(--primary), white text, padding 9px 18px, border-radius 6px

RIGHT CARD — "Existing Categories" (list)
5 category rows, each row:
  Left: [emoji 20px] + [Name in 13px] + [Hindi·Marathi in 11px muted below]
  Right: two action buttons (ti-edit, ti-trash)
  Row bg: var(--bg2), border-radius 6px, padding 8px 10px

Rows:
  🌾 Farmers       / किसान · शेतकरी
  🎓 Students      / छात्र · विद्यार्थी
  👩 Women         / महिला · महिला
  💼 Business Owners / व्यापारी · व्यापारी
  👴 Senior Citizens / वरिष्ठ नागरिक

---

### PAGE 4: USERS

#### Stats Row — 4 cards
Card 1: ti-users blue  | "18,430" | "Total Users"
Card 2: ti-user-check green | "16,204" | "Verified"
Card 3: ti-user-plus orange | "1,204" | "This Week"
Card 4: ti-user-off red | "22" | "Inactive"

#### Users Table (full width card)
Columns: User | State | Occupation | Registered | Status | Actions

Each User cell has:
  - Avatar circle (28x28, initials, colored bg)
  - Name (500 weight, 13px)
  - Email below (11px, muted)

Avatar colors:
  RK → primary-light bg, primary text
  PS → accent-light bg, accent text
  AJ → success-bg bg, success text

Rows:
  1. Ramesh Kumar / ramesh@gmail.com | Rajasthan | Farmer | 2 days ago | Verified (green)
  2. Priya Sharma / priya@gmail.com | Maharashtra | Student | 5 days ago | Verified (green)
  3. Amit Joshi / amit@gmail.com | UP | Business | 1 week ago | Unverified (gray)

Action buttons per row:
  ti-eye (view), ti-ban (deactivate — danger hover)

---

### PAGE 5: ELIGIBILITY RULES

Single full-width card: "Set Eligibility Rules for a Scheme"

Form fields:
  1. "Select Scheme" → select dropdown:
     Options: PM Kisan Samman Nidhi, NSP Scholarship

  2. Two-column row:
     - "Min Age" → number input, placeholder "18"
     - "Max Age" → number input, placeholder "60 (0 = no limit)"

  3. Two-column row:
     - "Min Annual Income (₹)" → number input, placeholder "0"
     - "Max Annual Income (₹)" → number input, placeholder "200000"

  4. Two-column row:
     - "Gender" → select: All / Male / Female / Other
     - "Caste Categories" → select: All / SC / ST / OBC / General

  5. "Applicable States (leave empty = All India)" → text input
     placeholder: "e.g. Maharashtra, Rajasthan, UP"

  6. "Occupations" → text input
     placeholder: "e.g. farmer, student, unemployed"

  7. Submit: "Save Eligibility Rules" with ti-device-floppy icon
     bg var(--primary), white, padding 9px 18px

---

### PAGE 6: NOTIFICATIONS

Two equal columns:

LEFT CARD — "Send Notification" (form)
Fields:
  1. "Send To" → select: All Users / All Farmers / All Students / Specific State
  2. "Title" → text input, placeholder "Scheme deadline reminder"
  3. "Message" → textarea (4 rows), placeholder "Dear citizen, Ladli Behna Yojana..."
  4. "Type" → select: Deadline Reminder / New Scheme Alert / System Update
  5. Submit: "Send Notification" with ti-send icon

RIGHT CARD — "Recent Notifications Sent"
Activity list, 3 items:
  🟠 "Deadline Reminder — Ladli Behna Yojana"
     subtext: "Sent to 4,230 women users · 2 hours ago"

  🟢 "New Scheme Alert — PM Vishwakarma Yojana"
     subtext: "Sent to 18,430 users · 1 day ago"

  🔵 "Deadline Reminder — NSP Scholarship 2024"
     subtext: "Sent to 2,100 student users · 3 days ago"

---

### PAGE 7: ANALYTICS

#### Stats Row — 4 cards
Card 1: "92,430" | "Total Eligibility Checks" | "+12% this month"
Card 2: "34,220" | "Schemes Saved by Users" | "+8% this month"
Card 3: "PM Kisan" | "Most Viewed Scheme" | "14,200 views"
Card 4: "Maharashtra" | "Top State by Users" | "3,420 users"

#### Top 5 Schemes Card (full width)
Title: "Top 5 Most Applied Schemes"
List rows (same style as category rows but with rank numbers):

  Rank | Scheme Name             | Bar (120px wide) | Count
  1    | PM Kisan Samman Nidhi   | 100% (accent)    | 14,200
  2    | NSP Post-Matric Scholarship | 78%          | 11,080
  3    | PM Mudra Yojana         | 60%              | 8,520
  4    | Ladli Behna Yojana      | 48%              | 6,840
  5    | Ayushman Bharat Yojana  | 35%              | 4,970

Rank numbers: 13px, weight 500, width 20px fixed
Bar: 4px height, border-radius 2px
Bar 1 fill: var(--accent) [orange, to distinguish from sidebar nav bars]
Bars 2-5 fill: var(--primary)

---

## REUSABLE COMPONENT SPECS

### Badge Component
```
.badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 100px;
  font-size: 11px;
  font-weight: 500;
}
```
Variants:
  .active    → bg: success-bg, color: success
  .inactive  → bg: bg2, color: text2
  .central   → bg: primary-light, color: primary
  .state     → bg: accent-light, color: accent
  .expiring  → bg: warning-bg, color: warning

### Form Input Styling
```
.form-input {
  padding: 8px 10px;
  border: 0.5px solid var(--border);
  border-radius: 6px;
  font-size: 13px;
  background: var(--bg);
  color: var(--text);
  outline: none;
  width: 100%;
}
.form-input:focus {
  border-color: var(--primary);
}
```

### Table Styling
- th: bg2 background, uppercase, 11px, muted, 0.05em letter-spacing, padding 8px 12px
- td: padding 10px 12px, border-bottom 0.5px var(--border)
- tr:hover td: background var(--bg2)
- Last row: no border-bottom
- Border-collapse: collapse

### Action Button (icon button in tables)
```
28x28px square
border-radius: 6px
border: 0.5px solid var(--border)
background: var(--bg)
color: var(--text2)
hover → background: var(--bg2)
danger hover → background: danger-bg, color: danger, border-color: danger
icon: 15px
```

### Card Component
```
background: var(--bg)
border: 0.5px solid var(--border)
border-radius: var(--border-radius-md)  [8px]
```
Card header: padding 14px 16px, border-bottom 0.5px, flex row, space-between
Card body: padding 16px (or 0 when contains a full table)

---

## JAVASCRIPT BEHAVIOR

### Page Switching Function
```javascript
function showPage(pageId, clickedElement) {
  // 1. Hide ALL pages
  document.querySelectorAll('.page').forEach(p => {
    p.style.display = 'none';
  });

  // 2. Show target page
  document.getElementById('page-' + pageId).style.display = 'block';

  // 3. Update active nav item
  document.querySelectorAll('.sb-item').forEach(i => {
    i.classList.remove('active');
  });
  clickedElement.classList.add('active');

  // 4. Update page title
  const titles = {
    overview: 'Overview',
    schemes: 'Schemes',
    categories: 'Categories',
    users: 'Users',
    eligibility: 'Eligibility Rules',
    notifications: 'Notifications',
    analytics: 'Analytics'
  };
  document.getElementById('page-title').textContent = titles[pageId];

  // 5. Update action button label
  const actions = {
    overview: 'Add Scheme',
    schemes: 'Add New Scheme',
    categories: 'Add Category',
    users: 'Add Admin',
    eligibility: 'Save Rules',
    notifications: 'Send Notification',
    analytics: 'Export Report'
  };
  document.getElementById('top-action').innerHTML =
    `<i class="ti ti-plus"></i> ${actions[pageId]}`;
}
```

### Nav Item onclick Attributes
Each sidebar nav item has:
  onclick="showPage('overview', this)"
  onclick="showPage('schemes', this)"
  onclick="showPage('categories', this)"
  onclick="showPage('users', this)"
  onclick="showPage('eligibility', this)"
  onclick="showPage('notifications', this)"
  onclick="showPage('analytics', this)"

### Default State on Load
- Page "overview" is visible (display:block), all others display:none
- Nav item "Overview" has class "active"

---

## SPACING & TYPOGRAPHY SYSTEM

Font: Plus Jakarta Sans, weights 400 and 500 only (never 600 or 700)

Font sizes:
  Page title: 15px / 500
  Card title: 13px / 500
  Body text: 13px / 400
  Table header: 11px / 500 / uppercase
  Muted/secondary: 11-12px / 400 / color text2
  Badge text: 11px / 500
  Stat number: 22px / 500
  Section label (sidebar): 10px / uppercase / letter-spacing 0.08em

Spacing scale:
  Content padding: 20px
  Card padding: 14px 16px (header), 16px (body)
  Grid gaps: 12px (stats), 16px (two-col)
  Table cell padding: 10px 12px (td), 8px 12px (th)
  Nav item padding: 9px 10px
  Form field gap: 12px between fields

---

## ACCESSIBILITY

- All decorative icons get aria-hidden="true"
- Page has a visually-hidden h1 or h2 describing the dashboard
- Interactive elements (buttons, inputs, selects) are focusable
- Color is not the only indicator (badges have text labels too)

---

## WHAT TO BUILD — SUMMARY CHECKLIST

[ ] Single HTML file with embedded CSS and JS
[ ] Tabler Icons CDN + Plus Jakarta Sans CDN loaded in <head>
[ ] CSS variables defined in :root
[ ] Fixed sidebar (220px) with brand, nav, footer sections
[ ] Fixed topbar (56px) with title, search, bell, action button
[ ] 7 page content divs — only one visible at a time
[ ] Page 1: Overview — 4 stats + alert + 2-col section + activity feed
[ ] Page 2: Schemes — filter row + data table with 5 rows
[ ] Page 3: Categories — add form + existing list (two columns)
[ ] Page 4: Users — 4 stats + users table with avatars
[ ] Page 5: Eligibility — single form with all eligibility fields
[ ] Page 6: Notifications — send form + activity feed (two columns)
[ ] Page 7: Analytics — 4 stats + top 5 schemes bar list
[ ] JavaScript showPage() function wired to all nav items
[ ] All badge variants styled correctly
[ ] Table hover states working
[ ] Action button danger hover on delete/ban buttons
[ ] Form inputs with focus border color
[ ] Responsive enough for 1280px+ screens
[ ] No external JS libraries. No build step. Opens directly in browser.

---

## IMPORTANT DESIGN RULES

1. NO gradients anywhere — flat solid colors only
2. NO drop shadows — borders only (0.5px)
3. NO dark backgrounds on page content area — only sidebar is dark
4. ALL borders are 0.5px (except featured card accent = 2px)
5. Border color is always var(--border) = rgba(0,0,0,0.08)
6. Sidebar background is ONLY var(--primary) = #1a3a6b
7. Page content background is var(--bg3) = #f1f5f9
8. Cards are white (var(--bg)) with 0.5px border
9. Text on colored backgrounds uses dark shade of SAME color family
10. Sentence case everywhere — never ALL CAPS in content (only th uppercase via CSS)
11. Icon size: 16px inline, 18px decorative max in nav, 24px max in stat boxes
12. Two font weights only: 400 (regular) and 500 (medium/bold)

---

*End of prompt. Use this to generate the exact Yojana Portal Admin Dashboard.*
