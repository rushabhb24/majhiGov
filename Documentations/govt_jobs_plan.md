# Implementation Plan: Government Jobs (Govt Jobs) Portal Integration

This document specifies the technical design, database schemas, API specs, and frontend architecture required to integrate a fully-featured **Government Jobs** listing and management module into **Yojana Portal (majhiGov)**.

---

## 👥 User & Business Value Proposition
For 100k+ citizens (students, farmers, unemployed youths), welfare schemes and government employment opportunities represent the two most critical services. By centralizing government jobs next to active welfare schemes in a unified portal, MajhiGov becomes an indispensable, one-stop platform.

* **Citizen Experience**: Seamless search, qualification matching, deadline alerts, fee structure transparency, and a direct application gateway.
* **Admin Experience**: Curation control, qualification checklists, and automatic translations into Hindi/Marathi via the existing translation API.
* **CTO/Technical Core**: High-efficiency connection pooling, clean REST endpoints, safe nullable database scanner mappings, and strict role-based access control.

---

## 🛠️ 1. Proposed Database Schema Changes

A new table `govt_jobs` will be introduced inside database migrations (`InitDB()` in `db.go`) to model job listings:

```sql
CREATE TABLE IF NOT EXISTS govt_jobs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    title_hi VARCHAR(255) NOT NULL DEFAULT '',
    title_mr VARCHAR(255) NOT NULL DEFAULT '',
    organization VARCHAR(255) NOT NULL,                  -- e.g., "MPSC", "UPSC", "Railway Board"
    department VARCHAR(255) NOT NULL,                    -- e.g., "Revenue Department", "IT Cell"
    vacancies INTEGER DEFAULT 0,
    education_qualification VARCHAR(255) NOT NULL,        -- e.g., "10th Pass", "12th Pass", "Graduate"
    experience_required VARCHAR(255) DEFAULT 'None',
    required_documents TEXT[] NOT NULL,                  -- e.g., {"Aadhaar Card", "Graduation Certificate"}
    application_start_date DATE NOT NULL,
    application_end_date DATE NOT NULL,
    official_website TEXT NOT NULL,
    apply_link TEXT NOT NULL,
    application_fee VARCHAR(255) NOT NULL,               -- e.g., "Gen: ₹394, Reserved: ₹294"
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexing for rapid citizen searches & filters
CREATE INDEX IF NOT EXISTS idx_govt_jobs_active_end ON govt_jobs(is_active, application_end_date);
CREATE INDEX IF NOT EXISTS idx_govt_jobs_qualification ON govt_jobs(education_qualification);
```

---

## 🔌 2. API Endpoint Architecture

### Public / Citizen Endpoints (Rate Limit: 50 req/min)
* **`GET /api/jobs`**: Lists all active government jobs. Supports pagination, search, and filtering by qualification and department.
* **`GET /api/jobs/:id`**: Fetches comprehensive details for a single job listing, including checklist items and bilingual translations.

### Admin Endpoints (Auth & Admin JWT Required)
* **`POST /api/admin/jobs`**: Inserts a new job listing with auto-translation triggers for title and description parameters.
* **`PUT /api/admin/jobs/:id`**: Updates an existing job's qualifications, vacancy details, fees, and links.
* **`DELETE /api/admin/jobs/:id`**: Soft-deletes a job by toggling its `is_active` parameter.

---

## 🎨 3. Frontend Component Design

### Main Citizen Portal Integration
* **Navigation Link**: Add a **"💼 Govt Jobs"** link in `Header.vue` right next to "Saved Schemes".
* **`JobsView.vue` (New Page)**:
  * **Interactive Search Bar & Filter Controls**: Filter by Organization (e.g. UPSC, MPSC), Department, and Education Qualification (e.g. 10th Pass, Graduate).
  * **Job Card List Layout**:
    * Show vacancies count badge, organization name, application fee, and a visual deadline indicator.
    * Highlight jobs closing soon (within 7 days) using an animated red warning badge `⚠️ Closing soon!`.
  * **Personalized AI/Smart Job Matching**:
    * If a user is logged in, automatically recommend jobs matching their profile's `education_level` inside a glassmorphic recommendations carousel!
  * **Job Details Modal**:
    * Overlays complete description, document requirements, and fees side-by-side with an action button that securely opens the external apply link in a new window.

### Admin Dashboard Integration
* **Sidebar Menu**: Add a **"💼 Manage Jobs"** tab in `AdminSidebar.vue`.
* **`JobsTab.vue` (New Component)**:
  * **Grid Statistics Overview**: Quick counters for Total Active Jobs, Total Vacancies, Expiring Jobs, and Deactivated listings.
  * **Tactile Job List Table**: Displays jobs, organizations, qualifications, vacancies, and application deadlines.
  * **CRUD Operations Modal**:
    * An admin form to create/edit jobs in English, which automatically triggers translations into Hindi and Marathi via the background translation handler.
    * Toggle switches to easily activate or suspend job postings.

---

## 🧪 4. Verification & Testing Plan

### 1. Automated System Validation
* Compile and test Go backend binaries:
  ```powershell
  cd backend; go build .\cmd\server
  ```
* Test and verify Vite production compiler bundles:
  ```powershell
  cd frontend; npm run build
  ```

### 2. Manual Acceptance Criteria
1. **Admin Job Creation**: Log in as an Administrator, navigate to "Manage Jobs", fill out the job creation form (e.g., "Post Officer"), save, and verify that Hindi/Marathi translations populate perfectly.
2. **Citizen Job Matching**: Log in as a citizen with "Graduate" profile qualifications, navigate to the Jobs tab, and confirm that only jobs requiring "Graduate" or below display inside the **"Personalized Job Recommendations for You"** matching slider.
3. **Deadline alerts**: Set a job's `application_end_date` to 5 days from today and verify the card displays a red `⚠️ Closing soon!` visual alarm.
4. **PWA Offline Availability**: Confirmed by cache registration inside the service worker `sw.js` for job endpoints.
