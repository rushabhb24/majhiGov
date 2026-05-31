# Yojana Portal (majhiGov) — API Specification & Postman Testing Guide

This documentation specifies the complete REST API interface for the Yojana Portal (majhiGov) platform. It covers public routes, secure citizen endpoints with at-rest Aadhaar encryption, and administrative interfaces. 

---

## 🚀 Postman Quick-Start Guide

### 1. Set Up Environment Variables
To avoid copying and pasting variables between requests, create a new **Postman Environment** and define the following variables:
* `baseUrl`: `http://localhost:8080`
* `token`: *(Leave blank initially; this will hold the JWT token returned during login)*

### 2. Configure Authentication
For all **Protected (Citizen)** and **Admin** routes:
1. Go to the request's **Authorization** tab.
2. Select **Type**: `Bearer Token`.
3. Enter `{{token}}` as the token value.

### 3. Automated Token Extraction (Optional but Recommended)
In Postman, you can automatically capture the token upon login. In the **Post-response** (or **Tests**) tab of your `/api/auth/login` request, paste the following snippet:
```javascript
const response = pm.response.json();
if (response.token) {
    pm.environment.set("token", response.token);
}
```

---

## 👥 1. Public API Endpoints (No Authentication Required)

### 📌 `GET /api/schemes`
* **Role/Purpose**: Fetches the list of all active welfare schemes from the database, grouped by their respective categories. Supports offline caching fallbacks.
* **Headers**:
  * `Accept: application/json`
* **Request Payload**: None
* **Sample Response (`200 OK`)**:
  ```json
  [
    {
      "id": 1,
      "title": "PM-Kisan Samman Nidhi Yojana",
      "title_hi": "पीएम-किसान सम्मान निधि योजना",
      "title_mr": "पीएम-किसान सन्मान निधी योजना",
      "description": "An initiative by the Government of India that provides up to ₹6,000 per year...",
      "description_hi": "भारत सरकार की एक पहल...",
      "description_mr": "अल्पभूधारक आणि सीमांत शेतकऱ्यांना...",
      "category_id": 1,
      "category_name": "Farmers",
      "government_level": "central",
      "state": null,
      "benefits": "₹6,000 per year in 3 installments",
      "application_start_date": "2026-05-20",
      "application_end_date": "2026-09-30",
      "official_website": "https://pmkisan.gov.in/",
      "apply_link": "https://pmkisan.gov.in/",
      "is_active": true,
      "created_at": "2026-05-30T10:00:00Z"
    }
  ]
  ```

---

### 📌 `GET /api/schemes/:id`
* **Role/Purpose**: Retrieves complete relational details for a single scheme, including eligibility thresholds, mandatory document checklists, and bilingual localized FAQs.
* **URL Parameter**: `id` (Integer) — The database primary key of the scheme.
* **Request Payload**: None
* **Sample Response (`200 OK`)**:
  ```json
  {
    "id": 1,
    "title": "PM-Kisan Samman Nidhi Yojana",
    "category_id": 1,
    "government_level": "central",
    "benefits": "₹6,000 per year",
    "eligibility": {
      "min_age": 18,
      "max_age": 100,
      "gender": "all",
      "caste_categories": ["General", "OBC", "SC", "ST"],
      "min_income": 0,
      "max_income": 300000,
      "occupations": ["Farmer"],
      "employee_types": ["Unemployed", "Self-Employed"],
      "education_levels": ["None", "Primary", "10th Pass", "Graduate"],
      "disability_required": false
    },
    "documents": [
      {
        "id": 1,
        "document_name": "Aadhaar Card",
        "document_name_hi": "आधार कार्ड",
        "document_name_mr": "आधार कार्ड",
        "is_mandatory": true
      }
    ],
    "faqs": [
      {
        "id": 1,
        "question": "Is bank account linkage mandatory?",
        "answer": "Yes, your bank account must be linked with Aadhaar.",
        "question_hi": "क्या बैंक खाता लिंक करना अनिवार्य है?",
        "answer_hi": "हाँ, डीबीटी क्रेडिट के लिए लिंक होना अनिवार्य है.",
        "question_mr": "बँक खाते आधारशी लिंक करणे बंधनकारक आहे का?",
        "answer_mr": "होय, थेट लाभ हस्तांतरण जमा होण्यासाठी बंधनकारक आहे."
      }
    ]
  }
  ```

---

### 📌 `POST /api/eligibility-check`
* **Role/Purpose**: Instantly evaluates a custom citizen demographic profile against all active schemes and returns a structured analysis specifying which schemes the citizen is eligible for and detailed mismatch reasons for others. Protects against DDoS with strict rate limits (50 req/min).
* **Headers**:
  * `Content-Type: application/json`
* **Request Body (JSON)**:
  ```json
  {
    "age": 25,
    "gender": "male",
    "state": "Maharashtra",
    "district": "Pune",
    "caste_category": "OBC",
    "annual_income": 120000,
    "occupation": "Farmer",
    "employee_type": "Self-Employed",
    "education_level": "12th Pass",
    "is_disabled": false
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "eligible": [
      {
        "id": 1,
        "title": "PM-Kisan Samman Nidhi Yojana",
        "government_level": "central",
        "benefits": "₹6,000 per year"
      }
    ],
    "ineligible": [
      {
        "id": 3,
        "title": "Lado Deviprasad Scheme (Mahila Unnati)",
        "reasons": [
          "Gender criteria mismatch (Requires 'female')"
        ]
      }
    ]
  }
  ```

---

### 📌 `POST /api/auth/register`
* **Role/Purpose**: Registers a new citizen user. Automatically encrypts the sensitive 12-digit Aadhaar Card number using AES-256 GCM before committing it to the database at-rest.
* **Headers**:
  * `Content-Type: application/json`
* **Request Body (JSON)**:
  ```json
  {
    "email": "citizen@gmail.com",
    "phone": "9876543210",
    "password": "securepassword123",
    "full_name": "Rajesh Tukaram Patil",
    "date_of_birth": "1992-06-15",
    "gender": "Male",
    "state": "Maharashtra",
    "district": "Kolhapur",
    "caste_category": "OBC",
    "annual_income": 140000.00,
    "occupation": "Farmer",
    "employee_type": "Self-Employed",
    "education_level": "10th Pass",
    "is_disabled": false,
    "aadhaar": "555566667777"
  }
  ```
* **Sample Response (`201 Created`)**:
  ```json
  {
    "success": true,
    "message": "Citizen registration successful! Security keys generated.",
    "user_id": 14
  }
  ```

---

### 📌 `POST /api/auth/login`
* **Role/Purpose**: Authenticates a user (citizen or admin) and returns a signed JWT token containing standard claims (`user_id`, `is_admin`). Protected by rate limits (5 req/min).
* **Headers**:
  * `Content-Type: application/json`
* **Request Body (JSON)**:
  ```json
  {
    "email": "citizen@gmail.com",
    "password": "securepassword123"
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Login successful",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ey...",
    "profile": {
      "id": 14,
      "user_id": 14,
      "full_name": "Rajesh Tukaram Patil",
      "email": "citizen@gmail.com",
      "phone": "9876543210",
      "is_admin": false
    }
  }
  ```

---

## 🔒 2. Protected Citizen API Endpoints (JWT Required)

* **Important**: Add header `Authorization: Bearer {{token}}` for all requests in this section.

### 📌 `GET /api/user/profile`
* **Role/Purpose**: Retrieves the authenticated user's profile card. Automatically decrypts their Aadhaar card number on the fly using AES-256 GCM.
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "profile": {
      "id": 14,
      "user_id": 14,
      "full_name": "Rajesh Tukaram Patil",
      "date_of_birth": "1992-06-15",
      "gender": "Male",
      "state": "Maharashtra",
      "district": "Kolhapur",
      "caste_category": "OBC",
      "annual_income": 140000,
      "occupation": "Farmer",
      "employee_type": "Self-Employed",
      "education_level": "10th Pass",
      "is_disabled": false,
      "aadhaar": "555566667777",
      "email": "citizen@gmail.com",
      "phone": "9876543210"
    }
  }
  ```

---

### 📌 `PUT /api/user/profile`
* **Role/Purpose**: Modifies citizen profile metadata and re-encrypts updated Aadhaar cards. Triggers auto-updates inside Pinia datastores and recalculates suggestions.
* **Request Body (JSON)**:
  ```json
  {
    "full_name": "Rajesh Patil",
    "date_of_birth": "1992-06-15",
    "gender": "Male",
    "state": "Maharashtra",
    "district": "Kolhapur",
    "caste_category": "OBC",
    "annual_income": 155000.00,
    "occupation": "Farmer",
    "employee_type": "Self-Employed",
    "education_level": "12th Pass",
    "is_disabled": false,
    "aadhaar": "555566667777"
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Demographic profile successfully updated!"
  }
  ```

---

### 📌 `POST /api/user/saved`
* **Role/Purpose**: Bookmarks or unbookmarks a scheme. Used to manage the user's saved wishlist drawer, generating warning badges if deadlines approach.
* **Request Body (JSON)**:
  ```json
  {
    "scheme_id": 2
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Scheme bookmarked successfully!"
  }
  ```

---

### 📌 `GET /api/user/saved`
* **Role/Purpose**: Returns all schemes bookmarked by the logged-in user. Supporting PWA offline read cache fallbacks.
* **Sample Response (`200 OK`)**:
  ```json
  [
    {
      "id": 2,
      "title": "Post Matric Scholarship Scheme",
      "benefits": "100% tuition fee waiver",
      "application_end_date": "2026-06-10"
    }
  ]
  ```

---

### 📌 `POST /api/user/apply`
* **Role/Purpose**: Submits a formal application for a welfare scheme.
* **Request Body (JSON)**:
  ```json
  {
    "scheme_id": 1
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Application submitted successfully! Government verification initiated."
  }
  ```

---

### 📌 `GET /api/user/applications`
* **Role/Purpose**: Fetches the list of all applications submitted by the logged-in citizen. Feeds the visual progress stepper.
* **Sample Response (`200 OK`)**:
  ```json
  [
    {
      "id": 24,
      "scheme_id": 1,
      "scheme_title": "PM-Kisan Samman Nidhi Yojana",
      "status": "pending",
      "applied_at": "2026-05-30T12:00:00Z",
      "notes": "Under document verification check."
    }
  ]
  ```

---

### 📌 `POST /api/translate`
* **Role/Purpose**: Interfaces with the translation engine. Translates English strings into localized Hindi & Marathi during administrative additions.
* **Request Body (JSON)**:
  ```json
  {
    "text": "Land Record Document (7/12 Extract)",
    "target_langs": ["hi", "mr"]
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "translations": {
      "hi": "भूमि रिकॉर्ड दस्तावेज (7/12 उतारा)",
      "mr": "जमीन सातबारा उतारा (7/12 उतारा)"
    }
  }
  ```

---

## 👑 3. Protected Administrator API Endpoints (Admin JWT Required)

* **Important**: Requires header `Authorization: Bearer {{token}}` for a user profile with `is_admin = true`.

### 📌 `GET /api/admin/analytics`
* **Role/Purpose**: Synthesizes dashboard metric aggregates (expiring schemes, active/pending metrics, recent activity feed streams, category distribution lists).
* **Sample Response (`200 OK`)**:
  ```json
  {
    "total_schemes": 42,
    "total_users": 1845,
    "total_applications": 210,
    "pending_applications": 14,
    "recent_activity": [
      {
        "type": "user",
        "text": "User 'Ramesh Patil' registered from Maharashtra",
        "time_ago": "2 minutes ago"
      }
    ]
  }
  ```

---

### 📌 `POST /api/admin/schemes`
* **Role/Purpose**: Creates a new scheme with its associated eligibility thresholds, document checklist, and FAQ blocks in a single relational transaction.
* **Request Body (JSON)**:
  ```json
  {
    "title": "PM Mudra Loan Scheme",
    "title_hi": "पीएम मुद्रा ऋण योजना",
    "title_mr": "पीएम मुद्रा कर्ज योजना",
    "description": "Financial assistance for business startups.",
    "description_hi": "स्टार्टअप के लिए वित्तीय सहायता।",
    "description_mr": "स्टार्टअपसाठी आर्थिक मदत.",
    "category_id": 5,
    "government_level": "central",
    "state": null,
    "benefits": "Up to ₹10 Lakhs business loan",
    "application_start_date": "2026-05-01",
    "application_end_date": "2026-12-31",
    "official_website": "https://mudra.org.in",
    "apply_link": "https://mudra.org.in",
    "is_active": true,
    "eligibility": {
      "min_age": 18,
      "max_age": 65,
      "gender": "all",
      "caste_categories": ["General", "OBC", "SC", "ST"],
      "min_income": 0,
      "max_income": 1000000,
      "states": [],
      "occupations": ["Business Owner"],
      "employee_types": ["Self-Employed"],
      "education_levels": ["10th Pass", "Graduate"],
      "disability_required": false
    },
    "documents": [
      { "document_name": "Aadhaar Card", "document_name_hi": "आधार कार्ड", "document_name_mr": "आधार कार्ड", "is_mandatory": true },
      { "document_name": "Business Plan Pitch", "document_name_hi": "व्यवसाय योजना", "document_name_mr": "व्यवसाय आराखडा", "is_mandatory": false }
    ],
    "faqs": [
      { "question": "What is the APY age limit?", "answer": "18 to 40 years.", "question_hi": "आयु सीमा क्या है?", "answer_hi": "18 से 40 वर्ष.", "question_mr": "वयोमर्यादा काय आहे?", "answer_mr": "१८ ते ४० वर्षे." }
    ]
  }
  ```
* **Sample Response (`210 Created`)**:
  ```json
  {
    "success": true,
    "scheme_id": 5,
    "message": "Scheme and eligibility parameters successfully created!"
  }
  ```

---

### 📌 `PUT /api/admin/schemes/:id`
* **Role/Purpose**: Re-writes scheme properties, FAQs, and eligibility criteria tables dynamically inside an ACID SQL transaction.
* **URL Parameter**: `id` (Integer)
* **Request Body**: *Same JSON format as POST /api/admin/schemes*
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Scheme, FAQ, documents, and eligibility criteria updated successfully!"
  }
  ```

---

### 📌 `DELETE /api/admin/schemes/:id`
* **Role/Purpose**: Performs a soft deactivation toggle on the scheme, switching `is_active` to prevent foreign key violations.
* **URL Parameter**: `id` (Integer)
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Scheme successfully deactivated!",
    "active": false
  }
  ```

---

### 📌 `GET /api/admin/applications`
* **Role/Purpose**: Fetches the global citizen application queue. Decrypts citizen Aadhaar card numbers on-the-fly and loads notes securely.
* **Sample Response (`200 OK`)**:
  ```json
  [
    {
      "id": 24,
      "user_id": 14,
      "full_name": "Rajesh Patil",
      "email": "citizen@gmail.com",
      "phone": "9876543210",
      "scheme_id": 1,
      "scheme_title": "PM-Kisan Samman Nidhi Yojana",
      "government_level": "central",
      "status": "pending",
      "applied_at": "2026-05-30T12:00:00Z",
      "notes": "Under document verification check.",
      "updated_at": "2026-05-30T12:00:00Z",
      "aadhaar": "555566667777"
    }
  ]
  ```

---

### 📌 `POST /api/admin/applications/status`
* **Role/Purpose**: Approves or Rejects a citizen's application. Automatically writes citizen system alert logs, pushes status changes, and fires an automated SMS dispatch to their mobile device.
* **Request Body (JSON)**:
  ```json
  {
    "application_id": 24,
    "status": "approved",
    "notes": "Land record check completed successfully. Verified."
  }
  ```
* **Sample Response (`200 OK`)**:
  ```json
  {
    "success": true,
    "message": "Application successfully approved!"
  }
  ```
* **Backend Console Dispatched SMS Log**:
  ```
  2026/05/30 23:45:00 [SMS GATEWAY] Sending SMS notification to phone +91 9876543210: "Dear Rajesh Patil, your application for 'PM-Kisan Samman Nidhi Yojana' has been APPROVED. Remarks: Land record check completed successfully. Verified."
  ```

---

### 📌 `POST /api/admin/notifications`
* **Role/Purpose**: Broadcasts system notifications to custom cohorts (Farmers, Students, specific States) for broad reach outreach.
* **Request Body (JSON)**:
  ```json
  {
    "send_to": "All Farmers",
    "state": "",
    "title": "PM Kisan Installment Broadcast",
    "message": "The 18th installment of PM Kisan has been successfully deposited in all verified accounts.",
    "type": "System Update"
  }
  ```
* **Sample Response (`201 Created`)**:
  ```json
  {
    "success": true,
    "message": "Notification broadcast successfully sent to 418 users!",
    "count": 418
  }
  ```
