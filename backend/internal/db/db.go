package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"yojana-portal/backend/internal/models"
)

var DB *sql.DB

// InitDB initializes connection to 'majhigov' and creates tables / seeds if needed
func InitDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")

	if host == "" { host = "localhost" }
	if port == "" { port = "5432" }
	if user == "" { user = "postgres" }
	if dbname == "" { dbname = "majhigov" }
	if sslmode == "" { sslmode = "disable" }

	// Connect to default 'postgres' database to check/create the target db
	connStrDefault := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		host, port, user, password, sslmode)
	
	dbDefault, err := sql.Open("postgres", connStrDefault)
	if err != nil {
		return fmt.Errorf("failed to open default postgres connection: %v", err)
	}
	defer dbDefault.Close()

	err = dbDefault.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping default postgres server: %v. Check password/connection.", err)
	}

	var exists bool
	queryCheckDB := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = dbDefault.QueryRow(queryCheckDB, dbname).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	if !exists {
		log.Printf("Database %s does not exist. Creating it...", dbname)
		_, err = dbDefault.Exec(fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("Database %s created successfully!", dbname)
	}

	// Now connect to 'majhigov'
	connStrApp := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	DB, err = sql.Open("postgres", connStrApp)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %v", dbname, err)
	}

	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping %s: %v", dbname, err)
	}

	// Configure database connection pooling limits to support high concurrency
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(30 * time.Minute)

	log.Printf("Successfully connected to database: %s with connection pooling limits set (MaxOpen: 100, MaxIdle: 10)", dbname)

	// Run Relational Migrations
	err = runMigrations()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	// Seed Relational Mock Data
	err = seedRelationalData()
	if err != nil {
		return fmt.Errorf("failed to seed data: %v", err)
	}

	// Seed Mock Applications
	err = seedMockApplications()
	if err != nil {
		log.Printf("Warning: Failed to seed mock applications: %v", err)
	}

	// Seed Mock Government Jobs
	err = seedMockJobs()
	if err != nil {
		log.Printf("Warning: Failed to seed mock government jobs: %v", err)
	}

	return nil
}

func runMigrations() error {
	queries := []string{
		// 1. Users Table
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			phone VARCHAR(20) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			is_verified BOOLEAN DEFAULT FALSE,
			is_admin BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 2. User Profiles Table
		`CREATE TABLE IF NOT EXISTS user_profiles (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			full_name VARCHAR(255) NOT NULL,
			date_of_birth DATE NOT NULL,
			gender VARCHAR(50) NOT NULL,
			state VARCHAR(100) NOT NULL,
			district VARCHAR(100) NOT NULL,
			caste_category VARCHAR(100) NOT NULL,
			annual_income NUMERIC(15,2) NOT NULL,
			occupation VARCHAR(100) NOT NULL,
			employee_type VARCHAR(100) NOT NULL,
			education_level VARCHAR(100) NOT NULL,
			is_disabled BOOLEAN DEFAULT FALSE,
			disability_type VARCHAR(100),
			avatar_url VARCHAR(500) DEFAULT '',
			aadhaar_encrypted TEXT DEFAULT '',
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 3. Scheme Categories Table
		`CREATE TABLE IF NOT EXISTS scheme_categories (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			name_hi VARCHAR(100) NOT NULL,
			name_mr VARCHAR(100) NOT NULL,
			icon VARCHAR(100) NOT NULL,
			description TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 4. Schemes Table
		`CREATE TABLE IF NOT EXISTS schemes (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			title_hi VARCHAR(255) NOT NULL,
			title_mr VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			description_hi TEXT NOT NULL,
			description_mr TEXT NOT NULL,
			category_id INTEGER REFERENCES scheme_categories(id) ON DELETE RESTRICT,
			government_level VARCHAR(50) NOT NULL, -- central/state
			state VARCHAR(100),                    -- NULL if central
			benefits TEXT NOT NULL,
			application_start_date DATE NOT NULL,
			application_end_date DATE NOT NULL,
			official_website TEXT NOT NULL,
			apply_link TEXT NOT NULL,
			is_active BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 5. Eligibility Criteria Table
		`CREATE TABLE IF NOT EXISTS eligibility_criteria (
			id SERIAL PRIMARY KEY,
			scheme_id INTEGER UNIQUE REFERENCES schemes(id) ON DELETE CASCADE,
			min_age INTEGER NOT NULL,
			max_age INTEGER NOT NULL,
			gender VARCHAR(50) NOT NULL,           -- all/male/female/other
			caste_categories TEXT[] NOT NULL,      -- Array type
			min_income NUMERIC(15,2) NOT NULL,
			max_income NUMERIC(15,2) NOT NULL,
			states TEXT[],                         -- Array type, NULL/empty = all India
			occupations TEXT[] NOT NULL,
			employee_types TEXT[] NOT NULL,
			education_levels TEXT[] NOT NULL,
			disability_required BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 6. Scheme Documents Table
		`CREATE TABLE IF NOT EXISTS scheme_documents (
			id SERIAL PRIMARY KEY,
			scheme_id INTEGER REFERENCES schemes(id) ON DELETE CASCADE,
			document_name VARCHAR(255) NOT NULL,
			document_name_hi VARCHAR(255) NOT NULL,
			document_name_mr VARCHAR(255) NOT NULL,
			is_mandatory BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 7. Scheme FAQs Table
		`CREATE TABLE IF NOT EXISTS scheme_faqs (
			id SERIAL PRIMARY KEY,
			scheme_id INTEGER REFERENCES schemes(id) ON DELETE CASCADE,
			question TEXT NOT NULL,
			answer TEXT NOT NULL,
			question_hi TEXT NOT NULL,
			answer_hi TEXT NOT NULL,
			question_mr TEXT NOT NULL,
			answer_mr TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 8. User Saved Schemes Table
		`CREATE TABLE IF NOT EXISTS user_saved_schemes (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			scheme_id INTEGER REFERENCES schemes(id) ON DELETE CASCADE,
			saved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 9. User Applied Schemes Table
		`CREATE TABLE IF NOT EXISTS user_applied_schemes (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			scheme_id INTEGER REFERENCES schemes(id) ON DELETE CASCADE,
			status VARCHAR(50) DEFAULT 'pending', -- pending, approved, rejected
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			notes TEXT,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 10. OTP Table
		`CREATE TABLE IF NOT EXISTS otps (
			id SERIAL PRIMARY KEY,
			contact VARCHAR(255) NOT NULL,
			otp_code VARCHAR(10) NOT NULL,
			expires_at TIMESTAMP NOT NULL,
			is_used BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// 11. Notifications Table
		`CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			message TEXT NOT NULL,
			type VARCHAR(50) NOT NULL,
			is_read BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,

		// Schema Alter Queries to cleanly update columns in case DB is already initialized
		`ALTER TABLE scheme_documents ADD COLUMN IF NOT EXISTS document_name_mr VARCHAR(255) DEFAULT '';`,
		`ALTER TABLE scheme_faqs ADD COLUMN IF NOT EXISTS question_mr TEXT DEFAULT '';`,
		`ALTER TABLE scheme_faqs ADD COLUMN IF NOT EXISTS answer_mr TEXT DEFAULT '';`,
		`ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500) DEFAULT '';`,
		`ALTER TABLE user_profiles ADD COLUMN IF NOT EXISTS aadhaar_encrypted TEXT DEFAULT '';`,
		`ALTER TABLE users ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`,

		// Backfill Marathi translations for Categories
		`UPDATE scheme_categories SET name_mr = 'शेतकरी' WHERE name = 'Farmers' AND (name_mr = '' OR name_mr IS NULL);`,
		`UPDATE scheme_categories SET name_mr = 'विद्यार्थी' WHERE name = 'Students' AND (name_mr = '' OR name_mr IS NULL);`,
		`UPDATE scheme_categories SET name_mr = 'महिला' WHERE name = 'Women' AND (name_mr = '' OR name_mr IS NULL);`,
		`UPDATE scheme_categories SET name_mr = 'ज्येष्ठ नागरिक' WHERE name = 'Senior Citizens' AND (name_mr = '' OR name_mr IS NULL);`,
		`UPDATE scheme_categories SET name_mr = 'व्यावसायिक' WHERE name = 'Business Owners' AND (name_mr = '' OR name_mr IS NULL);`,

		// Backfill Marathi translations for Schemes
		`UPDATE schemes SET title_mr = 'पीएम-किसान सन्मान निधी योजना', description_mr = 'अल्पभूधारक आणि सीमांत शेतकऱ्यांना तीन समान हप्त्यांमध्ये दरवर्षी ₹६,००० पर्यंत मदत देणारा भारत सरकारचा एक उपक्रम.' WHERE title = 'PM-Kisan Samman Nidhi Yojana' AND (title_mr = '' OR title_mr IS NULL);`,
		`UPDATE schemes SET title_mr = 'मॅट्रिक्युलेशन नंतरची शिष्यवृत्ती योजना', description_mr = 'अनुसूचित जाती, जमाती आणि इतर मागासवर्गीय विद्यार्थ्यांना माध्यमिक शिक्षणानंतरचे उच्च शिक्षण घेण्यासाठी सरकारकडून दिली जाणारी आर्थिक मदत.' WHERE title = 'Post Matric Scholarship Scheme' AND (title_mr = '' OR title_mr IS NULL);`,
		`UPDATE schemes SET title_mr = 'लाडो देवीप्रसाद योजना (महिला उन्नती)', description_mr = 'कमी उत्पन्न असणाऱ्या कुटुंबातील महिलांना स्वावलंबी बनवण्यासाठी मासिक आर्थिक सहाय्य आणि अनुदान देण्याचा राज्य सरकारचा उपक्रम.' WHERE title = 'Lado Deviprasad Scheme (Mahila Unnati)' AND (title_mr = '' OR title_mr IS NULL);`,
		`UPDATE schemes SET title_mr = 'अटल पेन्शन योजना (एपीवाय)', description_mr = 'असंघटित क्षेत्रातील कामगारांना वृद्धापकाळात आर्थिक सुरक्षितता मिळावी यासाठी भारत सरकारने सुरू केलेली पेन्शन योजना.' WHERE title = 'Atal Pension Yojana (APY)' AND (title_mr = '' OR title_mr IS NULL);`,

		// Backfill Marathi translations for Documents
		`UPDATE scheme_documents SET document_name_mr = 'आधार कार्ड' WHERE document_name = 'Aadhaar Card' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'जमीन सातबारा उतारा (7/12 उतारा)' WHERE document_name = 'Land Record Document (7/12 Extract)' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'बँक खाते पासबुक' WHERE document_name = 'Bank Account Passbook' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'जातीचा दाखला' WHERE document_name = 'Caste Certificate' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'उत्पन्नाचा दाखला' WHERE document_name = 'Income Certificate' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'मागील परीक्षेचे गुणपत्रक' WHERE document_name = 'Mark Sheet of Last Passed Exam' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'कॉलेज प्रवेश पावती' WHERE document_name = 'College Admission Receipt' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'अधिवास प्रमाणपत्र (डोमिसाईल)' WHERE document_name = 'State Domicile Certificate' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'उत्पन्नाचा दाखला (१.५ लाखांपेक्षा कमी)' WHERE document_name = 'Family Income Certificate (Family Income < 1.5L)' AND (document_name_mr = '' OR document_name_mr IS NULL);`,
		`UPDATE scheme_documents SET document_name_mr = 'बचत बँक खाते तपशील' WHERE document_name = 'Savings Bank Account Details' AND (document_name_mr = '' OR document_name_mr IS NULL);`,

		// Backfill Marathi translations for FAQs
		`UPDATE scheme_faqs SET question_mr = 'पीएम-किसान योजनेसाठी कोण पात्र नाही?', answer_mr = 'राज्य/केंद्र सरकारी कर्मचारी, आयकर भरणारे शेतकरी आणि संस्थात्मक जमीनधारक पात्र नाहीत.' WHERE question = 'Who is not eligible for PM-Kisan?' AND (question_mr = '' OR question_mr IS NULL);`,
		`UPDATE scheme_faqs SET question_mr = 'बँक खाते आधारशी लिंक करणे बंधनकारक आहे का?', answer_mr = 'होय, थेट लाभ हस्तांतरण (DBT) जमा होण्यासाठी बँक खाते आधारशी जोडणे बंधनकारक आहे.' WHERE question = 'Is bank account linkage mandatory?' AND (question_mr = '' OR question_mr IS NULL);`,
		`UPDATE scheme_faqs SET question_mr = 'खुल्या (General) प्रवर्गातील विद्यार्थी अर्ज करू शकतात का?', answer_mr = 'नाही, ही योजना केवळ अनुसूचित जाती (SC), अनुसूचित जमाती (ST) आणि इतर मागासवर्ग (OBC) विद्यार्थ्यांसाठी मर्यादित आहे. खुल्या वर्गातील विद्यार्थी इतर योजना तपासू शकतात.' WHERE question = 'Can General category students apply?' AND (question_mr = '' OR question_mr IS NULL);`,
		`UPDATE scheme_faqs SET question_mr = 'कमाल उत्पन्न मर्यादा काय आहे?', answer_mr = 'सर्व स्रोतांकडून मिळणारे वार्षिक कौटुंबिक उत्पन्न ₹२.५ लाखांपेक्षा जास्त नसावे.' WHERE question = 'What is the maximum income limit?' AND (question_mr = '' OR question_mr IS NULL);`,
		`UPDATE scheme_faqs SET question_mr = 'सरकारी सेवेत कार्यरत असलेल्या महिला पात्र आहेत का?', answer_mr = 'नाही, ज्या महिलांच्या कुटुंबातील सदस्य आयकर भरतात किंवा सरकारी नोकरीत आहेत त्या महिला या योजनेसाठी पात्र नाहीत.' WHERE question = 'Are working women in government service eligible?' AND (question_mr = '' OR question_mr IS NULL);`,
		`UPDATE scheme_faqs SET question_mr = 'एपीवाय (APY) मध्ये सामील होण्यासाठी वयोमर्यादा काय आहे?', answer_mr = 'तुम्ही १८ ते ४० वर्षे वयोगटा दरम्यान अटल पेन्शन योजनेमध्ये सामील होऊ शकता.' WHERE question = 'What is the age limit for joining APY?' AND (question_mr = '' OR question_mr IS NULL);`,
		`CREATE INDEX IF NOT EXISTS idx_user_applied_schemes_user_status ON user_applied_schemes(user_id, status);`,
		`CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);`,
		`CREATE INDEX IF NOT EXISTS idx_schemes_cat_active ON schemes(category_id, is_active);`,

		// 12. Government Jobs Table
		`CREATE TABLE IF NOT EXISTS government_jobs (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			title_hi VARCHAR(255) NOT NULL,
			title_mr VARCHAR(255) NOT NULL,
			organization VARCHAR(255) NOT NULL,
			organization_hi VARCHAR(255) NOT NULL,
			organization_mr VARCHAR(255) NOT NULL,
			description TEXT NOT NULL,
			description_hi TEXT NOT NULL,
			description_mr TEXT NOT NULL,
			education_qualification VARCHAR(255) NOT NULL,
			documents_required TEXT[] NOT NULL,
			min_age INTEGER DEFAULT 18,
			max_age INTEGER DEFAULT 45,
			last_date DATE NOT NULL,
			apply_link TEXT NOT NULL,
			general_fee NUMERIC(15,2) DEFAULT 0.00,
			obc_fee NUMERIC(15,2) DEFAULT 0.00,
			sc_st_fee NUMERIC(15,2) DEFAULT 0.00,
			women_fee NUMERIC(15,2) DEFAULT 0.00,
			is_active BOOLEAN DEFAULT TRUE,
			clicks_count INTEGER DEFAULT 0,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE INDEX IF NOT EXISTS idx_jobs_active_date ON government_jobs(is_active, last_date);`,
	}

	for idx, query := range queries {
		_, err := DB.Exec(query)
		if err != nil {
			return fmt.Errorf("failed migration table index %d: %v", idx, err)
		}
	}

	log.Println("All 11 relational database tables are ready!")
	return nil
}

func seedRelationalData() error {
	// Seed default Super Admin if not present
	var adminExists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = 'admin@gov.in')").Scan(&adminExists)
	if err != nil {
		return fmt.Errorf("failed to check if admin exists: %v", err)
	}
	if !adminExists {
		log.Println("Seeding default Super Admin user...")
		hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash admin password: %v", err)
		}
		var adminID int
		err = DB.QueryRow(`
			INSERT INTO users (email, phone, password_hash, is_verified, is_admin)
			VALUES ('admin@gov.in', '9999999999', $1, true, true) RETURNING id`, string(hash)).Scan(&adminID)
		if err != nil {
			return fmt.Errorf("failed to insert seeded admin user: %v", err)
		}
		_, err = DB.Exec(`
			INSERT INTO user_profiles (
				user_id, full_name, date_of_birth, gender, state, district,
				caste_category, annual_income, occupation, employee_type,
				education_level, is_disabled
			) VALUES ($1, 'Super Admin', '1990-01-01', 'Male', 'Maharashtra', 'Mumbai',
			          'General', 0.00, 'Other', 'Government Employee', 'Graduate', false)`, adminID)
		if err != nil {
			return fmt.Errorf("failed to insert seeded admin profile: %v", err)
		}
		log.Println("Default Super Admin user (admin@gov.in / admin123) seeded successfully!")
	}

	// 1. Check if Categories already exist
	var catCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM scheme_categories").Scan(&catCount)
	if err != nil {
		return err
	}

	if catCount > 0 {
		log.Println("Database already contains seeded categories/schemes. Skipping seeding.")
		return nil
	}

	log.Println("Seeding scheme categories...")
	
	// Seed Categories
	categories := []struct {
		Name        string
		NameHi      string
		NameMr      string
		Icon        string
		Description string
	}{
		{"Farmers", "किसान", "शेतकरी", "tractor", "Agriculture support and farmer benefits"},
		{"Students", "छात्र", "विद्यार्थी", "graduation-cap", "Scholarships and educational loans"},
		{"Women", "महिलाएं", "महिला", "user-female", "Welfare and entrepreneurship schemes for women"},
		{"Senior Citizens", "वरिष्ठ नागरिक", "ज्येष्ठ नागरिक", "heart", "Pensions and health support schemes"},
		{"Business Owners", "व्यवसायी", "व्यावसायिक", "briefcase", "Subsidies, grants, and startup loans"},
	}

	var categoryIds = make(map[string]int)
	for _, c := range categories {
		var id int
		query := `
		INSERT INTO scheme_categories (name, name_hi, name_mr, icon, description)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
		err = DB.QueryRow(query, c.Name, c.NameHi, c.NameMr, c.Icon, c.Description).Scan(&id)
		if err != nil {
			return fmt.Errorf("failed to insert category %s: %v", c.Name, err)
		}
		categoryIds[c.Name] = id
	}

	log.Println("Categories seeded! Seeding schemes, eligibility criteria, documents, and FAQs...")

	// Seed Schemes Relational
	type fullSchemeSeed struct {
		Scheme      models.Scheme
		Eligibility models.EligibilityCriteria
		Docs        []models.SchemeDocument
		FAQs        []models.SchemeFAQ
	}

	seedData := []fullSchemeSeed{
		{
			Scheme: models.Scheme{
				Title:           "PM-Kisan Samman Nidhi Yojana",
				TitleHi:         "पीएम-किसान सम्मान निधि योजना",
				TitleMr:         "पीएम-किसान सन्मान निधी योजना",
				Description:     "An initiative by the Government of India that provides up to ₹6,000 per year in three equal installments to small and marginal farmers.",
				DescriptionHi:   "भारत सरकार की एक पहल जो छोटे और सीमांत किसानों को तीन समान किश्तों में प्रति वर्ष ₹6,000 तक प्रदान करती है।",
				DescriptionMr:   "अल्पभूधारक आणि सीमांत शेतकऱ्यांना तीन समान हप्त्यांमध्ये दरवर्षी ₹६,००० पर्यंत मदत देणारा भारत सरकारचा एक उपक्रम.",
				CategoryID:      categoryIds["Farmers"],
				GovernmentLevel: "central",
				Benefits:        "₹6,000 per year in 3 installments of ₹2,000 (₹6,000 प्रति वर्ष 3 किश्तों में / दरवर्षी ₹६,००० तीन हप्त्यात)",
				OfficialWebsite: "https://pmkisan.gov.in/",
				ApplyLink:       "https://pmkisan.gov.in/",
			},
			Eligibility: models.EligibilityCriteria{
				MinAge:             18,
				MaxAge:             100,
				Gender:             "all",
				CasteCategories:    pq.StringArray{"General", "OBC", "SC", "ST"},
				MinIncome:          0,
				MaxIncome:          300000,
				States:             pq.StringArray{}, // All India
				Occupations:        pq.StringArray{"Farmer"},
				EmployeeTypes:      pq.StringArray{"Unemployed", "Self-Employed"},
				EducationLevels:    pq.StringArray{"None", "Primary", "10th Pass", "12th Pass", "Graduate", "Post Graduate"},
				DisabilityRequired: false,
			},
			Docs: []models.SchemeDocument{
				{DocumentName: "Aadhaar Card", DocumentNameHi: "आधार कार्ड", DocumentNameMr: "आधार कार्ड", IsMandatory: true},
				{DocumentName: "Land Record Document (7/12 Extract)", DocumentNameHi: "भूमि रिकॉर्ड दस्तावेज (7/12 उतारा)", DocumentNameMr: "जमीन सातबारा उतारा (7/12 उतारा)", IsMandatory: true},
				{DocumentName: "Bank Account Passbook", DocumentNameHi: "बैंक खाता पासबुक", DocumentNameMr: "बँक खाते पासबुक", IsMandatory: true},
			},
			FAQs: []models.SchemeFAQ{
				{
					Question:   "Who is not eligible for PM-Kisan?",
					Answer:     "State/Central government employees, income taxpayers, and institutional landholders are not eligible.",
					QuestionHi: "पीएम-किसान के लिए कौन पात्र नहीं है?",
					AnswerHi:   "राज्य/केंद्र सरकार के कर्मचारी, आयकर दाता और संस्थागत भूमि धारक पात्र नहीं हैं.",
					QuestionMr: "पीएम-किसान योजनेसाठी कोण पात्र नाही?",
					AnswerMr:   "राज्य/केंद्र सरकारी कर्मचारी, आयकर भरणारे शेतकरी आणि संस्थात्मक जमीनधारक पात्र नाहीत.",
				},
				{
					Question:   "Is bank account linkage mandatory?",
					Answer:     "Yes, your bank account must be linked with Aadhaar for DBT credit.",
					QuestionHi: "क्या बैंक खाता लिंक करना अनिवार्य है?",
					AnswerHi:   "हाँ, डीबीटी क्रेडिट के लिए आपका बैंक खाता आधार से लिंक होना अनिवार्य है.",
					QuestionMr: "बँक खाते आधारशी लिंक करणे बंधनकारक आहे का?",
					AnswerMr:   "होय, थेट लाभ हस्तांतरण (DBT) जमा होण्यासाठी बँक खाते आधारशी जोडणे बंधनकारक आहे.",
				},
			},
		},
		{
			Scheme: models.Scheme{
				Title:           "Post Matric Scholarship Scheme",
				TitleHi:         "मैट्रिकोत्तर छात्रवृत्ति योजना",
				TitleMr:         "मॅट्रिक्युलेशन नंतरची शिष्यवृत्ती योजना",
				Description:     "Financial assistance provided by the government to students belonging to scheduled castes, tribes, and backward classes to pursue post-secondary education.",
				DescriptionHi:   "अनुसूचित जाति, जनजाति और पिछड़े वर्ग के छात्रों को उच्च शिक्षा प्राप्त करने के लिए सरकार द्वारा प्रदान की जाने वाली वित्तीय सहायता।",
				DescriptionMr:   "अनुसूचित जाती, जमाती आणि इतर मागासवर्गीय विद्यार्थ्यांना माध्यमिक शिक्षणानंतरचे उच्च शिक्षण घेण्यासाठी सरकारकडून दिली जाणारी आर्थिक मदत.",
				CategoryID:      categoryIds["Students"],
				GovernmentLevel: "central",
				Benefits:        "100% tuition fee waiver and monthly maintenance allowance (100% शिक्षण शुल्क माफी और मासिक भत्ता / १००% शिक्षण शुल्क माफी आणि मासिक भत्ता)",
				OfficialWebsite: "https://scholarships.gov.in/",
				ApplyLink:       "https://scholarships.gov.in/",
			},
			Eligibility: models.EligibilityCriteria{
				MinAge:             15,
				MaxAge:             30,
				Gender:             "all",
				CasteCategories:    pq.StringArray{"SC", "ST", "OBC"},
				MinIncome:          0,
				MaxIncome:          250000,
				States:             pq.StringArray{}, // All India
				Occupations:        pq.StringArray{"Student"},
				EmployeeTypes:      pq.StringArray{"Unemployed"},
				EducationLevels:    pq.StringArray{"10th Pass", "12th Pass", "Graduate"},
				DisabilityRequired: false,
			},
			Docs: []models.SchemeDocument{
				{DocumentName: "Caste Certificate", DocumentNameHi: "जाति प्रमाणपत्र", DocumentNameMr: "जातीचा दाखला", IsMandatory: true},
				{DocumentName: "Income Certificate", DocumentNameHi: "आय प्रमाणपत्र", DocumentNameMr: "उत्पन्नाचा दाखला", IsMandatory: true},
				{DocumentName: "Mark Sheet of Last Passed Exam", DocumentNameHi: "पिछली परीक्षा की मार्कशीट", DocumentNameMr: "मागील परीक्षेचे गुणपत्रक", IsMandatory: true},
				{DocumentName: "College Admission Receipt", DocumentNameHi: "कॉलेज प्रवेश रसीद", DocumentNameMr: "कॉलेज प्रवेश पावती", IsMandatory: true},
			},
			FAQs: []models.SchemeFAQ{
				{
					Question:   "Can General category students apply?",
					Answer:     "No, this specific scheme is restricted to SC, ST, and OBC students. General students can check other NSP schemes.",
					QuestionHi: "क्या सामान्य श्रेणी के छात्र आवेदन कर सकते हैं?",
					AnswerHi:   "नहीं, यह योजना केवल एससी, एसटी और ओबीसी छात्रों के लिए है. सामान्य वर्ग के छात्र अन्य एनएसपी योजनाओं की जांच कर सकते हैं.",
					QuestionMr: "खुल्या (General) प्रवर्गातील विद्यार्थी अर्ज करू शकतात का?",
					AnswerMr:   "नाही, ही योजना केवळ अनुसूचित जाती (SC), अनुसूचित जमाती (ST) आणि इतर मागासवर्ग (OBC) विद्यार्थ्यांसाठी मर्यादित आहे. खुल्या वर्गातील विद्यार्थी इतर योजना तपासू शकतात.",
				},
				{
					Question:   "What is the maximum income limit?",
					Answer:     "The family annual income from all sources must not exceed ₹2.5 Lakhs.",
					QuestionHi: "अधिकतम आय सीमा क्या है?",
					AnswerHi:   "सभी स्रोतों से पारिवारिक वार्षिक आय ₹2.5 लाख से अधिक नहीं होनी चाहिए.",
					QuestionMr: "कमाल उत्पन्न मर्यादा काय आहे?",
					AnswerMr:   "सर्व स्रोतांकडून मिळणारे वार्षिक कौटुंबिक उत्पन्न ₹२.५ लाखांपेक्षा जास्त नसावे.",
				},
			},
		},
		{
			Scheme: models.Scheme{
				Title:           "Lado Deviprasad Scheme (Mahila Unnati)",
				TitleHi:         "लाडो देवीप्रसाद योजना (महिला उन्नति)",
				TitleMr:         "लाडो देवीप्रसाद योजना (महिला उन्नती)",
				Description:     "A state-sponsored initiative aimed at providing monthly financial support and micro-grants to women from low-income families to foster self-reliance.",
				DescriptionHi:   "कम आय वाले परिवारों की महिलाओं को आत्मनिर्भर बनाने के लिए मासिक वित्तीय सहायता और सूक्ष्म अनुदान प्रदान करने के उद्देश्य से एक राज्य प्रायोजित पहल।",
				DescriptionMr:   "कमी उत्पन्न असणाऱ्या कुटुंबातील महिलांना स्वावलंबी बनवण्यासाठी मासिक आर्थिक सहाय्य आणि अनुदान देण्याचा राज्य सरकारचा उपक्रम.",
				CategoryID:      categoryIds["Women"],
				GovernmentLevel: "state",
				State:           stringPtr("Maharashtra"),
				Benefits:        "₹1,500 monthly transfer and up to ₹25,000 interest-free business grants (₹1,500 मासिक सहायता और ₹25,000 तक व्यवसाय अनुदान / दरमहा ₹१,५०० आणि ₹२५,००० पर्यंत बिनव्याजी व्यवसाय अनुदान)",
				OfficialWebsite: "https://wcd.gov.in/",
				ApplyLink:       "https://wcd.gov.in/",
			},
			Eligibility: models.EligibilityCriteria{
				MinAge:             18,
				MaxAge:             60,
				Gender:             "female",
				CasteCategories:    pq.StringArray{"General", "OBC", "SC", "ST"},
				MinIncome:          0,
				MaxIncome:          150000,
				States:             pq.StringArray{"Maharashtra", "Madhya Pradesh", "Gujarat"},
				Occupations:        pq.StringArray{"Farmer", "Student", "Business Owner", "Unemployed", "Self-Employed", "Other"},
				EmployeeTypes:      pq.StringArray{"Unemployed", "Self-Employed"},
				EducationLevels:    pq.StringArray{"None", "Primary", "10th Pass", "12th Pass", "Graduate", "Post Graduate"},
				DisabilityRequired: false,
			},
			Docs: []models.SchemeDocument{
				{DocumentName: "Aadhaar Card", DocumentNameHi: "आधार कार्ड", DocumentNameMr: "आधार कार्ड", IsMandatory: true},
				{DocumentName: "State Domicile Certificate", DocumentNameHi: "मूल निवासी प्रमाणपत्र", DocumentNameMr: "अधिवास प्रमाणपत्र (डोमिसाईल)", IsMandatory: true},
				{DocumentName: "Family Income Certificate (Family Income < 1.5L)", DocumentNameHi: "आय प्रमाणपत्र (1.5 लाख से कम)", DocumentNameMr: "उत्पन्नाचा दाखला (१.५ लाखांपेक्षा कमी)", IsMandatory: true},
			},
			FAQs: []models.SchemeFAQ{
				{
					Question:   "Are working women in government service eligible?",
					Answer:     "No, women whose family members pay income tax or are employed in government jobs are not eligible.",
					QuestionHi: "क्या सरकारी सेवा में कार्यरत महिलाएं पात्र हैं?",
					AnswerHi:   "नहीं, जिन महिलाओं के परिवार के सदस्य आयकर देते हैं या सरकारी नौकरियों में कार्यरत हैं, वे पात्र नहीं हैं.",
					QuestionMr: "सरकारी सेवेत कार्यरत असलेल्या महिला पात्र आहेत का?",
					AnswerMr:   "नाही, ज्या महिलांच्या कुटुंबातील सदस्य आयकर भरतात किंवा सरकारी नोकरीत आहेत त्या महिला या योजनेसाठी पात्र नाहीत.",
				},
			},
		},
		{
			Scheme: models.Scheme{
				Title:           "Atal Pension Yojana (APY)",
				TitleHi:         "अटल पेंशन योजना (एपीवाई)",
				TitleMr:         "अटल पेन्शन योजना (एपीवाय)",
				Description:     "A government-backed pension scheme in India, primarily targeted at the unorganized sector to provide financial security in old age.",
				DescriptionHi:   "भारत में एक सरकार समर्थित पेंशन योजना, मुख्य रूप से बुढ़ापे में वित्तीय सुरक्षा प्रदान करने के लिए असंगठित क्षेत्र को लक्षित करती है।",
				DescriptionMr:   "असंघटित क्षेत्रातील कामगारांना वृद्धापकाळात आर्थिक सुरक्षितता मिळावी यासाठी भारत सरकारने सुरू केलेली पेन्शन योजना.",
				CategoryID:      categoryIds["Senior Citizens"],
				GovernmentLevel: "central",
				Benefits:        "Guaranteed monthly pension of ₹1,000 to ₹5,000 after 60 years (₹1,000 से ₹5,000 की सुनिश्चित मासिक पेंशन / ६० वर्षांनंतर दरमहा ₹१,००० ते ₹५,००० पेन्शनची हमी)",
				OfficialWebsite: "https://www.npscra.nsdl.co.in/",
				ApplyLink:       "https://www.npscra.nsdl.co.in/",
			},
			Eligibility: models.EligibilityCriteria{
				MinAge:             18,
				MaxAge:             40,
				Gender:             "all",
				CasteCategories:    pq.StringArray{"General", "OBC", "SC", "ST"},
				MinIncome:          0,
				MaxIncome:          500000,
				States:             pq.StringArray{}, // All India
				Occupations:        pq.StringArray{"Farmer", "Student", "Business Owner", "Unemployed", "Self-Employed", "Other"},
				EmployeeTypes:      pq.StringArray{"Private", "Unemployed", "Self-Employed"},
				EducationLevels:    pq.StringArray{"None", "Primary", "10th Pass", "12th Pass", "Graduate", "Post Graduate"},
				DisabilityRequired: false,
			},
			Docs: []models.SchemeDocument{
				{DocumentName: "Aadhaar Card", DocumentNameHi: "आधार कार्ड", DocumentNameMr: "आधार कार्ड", IsMandatory: true},
				{DocumentName: "Savings Bank Account Details", DocumentNameHi: "बचत बैंक खाता विवरण", DocumentNameMr: "बचत बँक खाते तपशील", IsMandatory: true},
			},
			FAQs: []models.SchemeFAQ{
				{
					Question:   "What is the age limit for joining APY?",
					Answer:     "You can join APY between the age of 18 and 40 years.",
					QuestionHi: "एपीवाई में शामिल होने की आयु सीमा क्या है?",
					AnswerHi:   "आप 18 से 40 वर्ष की आयु के बीच अटल पेंशन योजना में शामिल हो सकते हैं.",
					QuestionMr: "एपीवाय (APY) मध्ये सामील होण्यासाठी वयोमर्यादा काय आहे?",
					AnswerMr:   "तुम्ही १८ ते ४० वर्षे वयोगटा दरम्यान अटल पेन्शन योजनेमध्ये सामील होऊ शकता.",
				},
			},
		},
	}

	for _, data := range seedData {
		var schemeId int
		// Insert Scheme
		queryScheme := `
		INSERT INTO schemes (title, title_hi, title_mr, description, description_hi, description_mr, category_id, government_level, state, benefits, application_start_date, application_end_date, official_website, apply_link)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, CURRENT_DATE - INTERVAL '10 days', CURRENT_DATE + INTERVAL '120 days', $11, $12)
		RETURNING id`
		err = DB.QueryRow(queryScheme, 
			data.Scheme.Title, data.Scheme.TitleHi, data.Scheme.TitleMr,
			data.Scheme.Description, data.Scheme.DescriptionHi, data.Scheme.DescriptionMr,
			data.Scheme.CategoryID, data.Scheme.GovernmentLevel, data.Scheme.State,
			data.Scheme.Benefits, data.Scheme.OfficialWebsite, data.Scheme.ApplyLink,
		).Scan(&schemeId)
		if err != nil {
			return fmt.Errorf("failed seeding scheme %s: %v", data.Scheme.Title, err)
		}

		// Insert Eligibility
		queryElig := `
		INSERT INTO eligibility_criteria (scheme_id, min_age, max_age, gender, caste_categories, min_income, max_income, states, occupations, employee_types, education_levels, disability_required)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
		_, err = DB.Exec(queryElig,
			schemeId, data.Eligibility.MinAge, data.Eligibility.MaxAge, data.Eligibility.Gender,
			data.Eligibility.CasteCategories, data.Eligibility.MinIncome, data.Eligibility.MaxIncome,
			data.Eligibility.States, data.Eligibility.Occupations, data.Eligibility.EmployeeTypes,
			data.Eligibility.EducationLevels, data.Eligibility.DisabilityRequired,
		)
		if err != nil {
			return fmt.Errorf("failed seeding eligibility for scheme id %d: %v", schemeId, err)
		}

		// Insert Documents
		for _, doc := range data.Docs {
			queryDoc := `
			INSERT INTO scheme_documents (scheme_id, document_name, document_name_hi, document_name_mr, is_mandatory)
			VALUES ($1, $2, $3, $4, $5)`
			_, err = DB.Exec(queryDoc, schemeId, doc.DocumentName, doc.DocumentNameHi, doc.DocumentNameMr, doc.IsMandatory)
			if err != nil {
				return fmt.Errorf("failed seeding document for scheme id %d: %v", schemeId, err)
			}
		}

		// Insert FAQs
		for _, faq := range data.FAQs {
			queryFAQ := `
			INSERT INTO scheme_faqs (scheme_id, question, answer, question_hi, answer_hi, question_mr, answer_mr)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`
			_, err = DB.Exec(queryFAQ, schemeId, faq.Question, faq.Answer, faq.QuestionHi, faq.AnswerHi, faq.QuestionMr, faq.AnswerMr)
			if err != nil {
				return fmt.Errorf("failed seeding FAQ for scheme id %d: %v", schemeId, err)
			}
		}
	}

	log.Println("Seeded relational categories, schemes, criteria, documents, and FAQs successfully!")
	return nil
}

func stringPtr(s string) *string {
	return &s
}

func seedMockApplications() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE is_admin = false").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	log.Println("Seeding mock citizen users and applications...")

	// 1. Get first 3 schemes
	var schemeIds []int
	rows, err := DB.Query("SELECT id FROM schemes LIMIT 3")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err == nil {
			schemeIds = append(schemeIds, id)
		}
	}

	if len(schemeIds) == 0 {
		log.Println("No schemes found to seed applications for. Skipping.")
		return nil
	}

	// 2. Create mock citizens
	mockUsers := []struct {
		Email      string
		Phone      string
		Name       string
		Occupation string
		State      string
		SchemeIdx  int
		Status     string
		Notes      string
		Aadhaar    string
	}{
		{"ramesh@gmail.com", "9876543211", "Ramesh Kumar", "Farmer", "Rajasthan", 0, "approved", "Eligible farmer with verified Aadhaar and land record parameters.", "987654321098"},
		{"priya@gmail.com", "9876543212", "Priya Sharma", "Student", "Maharashtra", 1, "pending", "Undergraduate student applying for Post Matric Scholarship scheme.", "555566667777"},
		{"amit@gmail.com", "9876543213", "Amit Joshi", "Business", "UP", 2, "rejected", "Business owner applying for subsidy. Annual income exceeds threshold limit.", "444455556666"},
	}

	hash, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	for _, mu := range mockUsers {
		var uid int
		// Insert User
		err = DB.QueryRow(`
			INSERT INTO users (email, phone, password_hash, is_verified, is_admin)
			VALUES ($1, $2, $3, true, false)
			ON CONFLICT (email) DO UPDATE SET email=EXCLUDED.email RETURNING id`,
			mu.Email, mu.Phone, string(hash)).Scan(&uid)
		if err != nil {
			// If already exists, query the ID
			err = DB.QueryRow("SELECT id FROM users WHERE email = $1", mu.Email).Scan(&uid)
			if err != nil {
				continue
			}
		}

		// Encrypt mock Aadhaar card
		aadhaarEncrypted, _ := Encrypt(mu.Aadhaar)

		// Insert Profile
		var profileExists bool
		_ = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id = $1)", uid).Scan(&profileExists)
		if !profileExists {
			_, _ = DB.Exec(`
				INSERT INTO user_profiles (
					user_id, full_name, date_of_birth, gender, state, district,
					caste_category, annual_income, occupation, employee_type,
					education_level, is_disabled, aadhaar_encrypted
				) VALUES ($1, $2, '1995-05-15', 'Male', $3, 'District Office',
				          'OBC', 120000.00, $4, 'Self-Employed', 'Graduate', false, $5)`,
				uid, mu.Name, mu.State, mu.Occupation, aadhaarEncrypted)
		}

		// Map Scheme
		schemeId := schemeIds[mu.SchemeIdx % len(schemeIds)]

		// Insert Application
		_, err = DB.Exec(`
			INSERT INTO user_applied_schemes (user_id, scheme_id, status, notes)
			VALUES ($1, $2, $3, $4)`,
			uid, schemeId, mu.Status, mu.Notes)
		if err != nil {
			log.Printf("Failed seeding application for %s: %v", mu.Email, err)
		}
	}

	log.Println("Mock citizen users and applications successfully seeded!")
	return nil
}

// seedMockJobs inserts initial government job advertisements
func seedMockJobs() error {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM government_jobs").Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		log.Println("Database already contains seeded government jobs. Skipping.")
		return nil
	}

	log.Println("Seeding mock government jobs...")

	jobs := []struct {
		Title                  string
		TitleHi                string
		TitleMr                string
		Org                    string
		OrgHi                  string
		OrgMr                  string
		Desc                   string
		DescHi                 string
		DescMr                 string
		EducationQualification string
		Docs                   []string
		MinAge                 int
		MaxAge                 int
		LastDate               string
		Link                   string
		GenFee                 float64
		ObcFee                 float64
		ScStFee                float64
		WomenFee               float64
	}{
		{
			Title: "Civil Services Examination (CSE) 2026",
			TitleHi: "सिविल सेवा परीक्षा (CSE) 2026",
			TitleMr: "नागरी सेवा परीक्षा (CSE) 2026",
			Org: "Union Public Service Commission (UPSC)",
			OrgHi: "संघ लोक सेवा आयोग (UPSC)",
			OrgMr: "केंद्रीय लोकसेवा आयोग (UPSC)",
			Desc: "Apply for IAS, IPS, IFS, and other Group A/B central services posts. Selection through Prelims, Mains, and Personality Test.",
			DescHi: "आईएएस, आईपीएस, आईएफएस और अन्य ग्रुप ए/बी केंद्रीय सेवा पदों के लिए आवेदन करें। प्रारंभिक, मुख्य और साक्षात्कार के माध्यम से चयन।",
			DescMr: "IAS, IPS, IFS आणि इतर गट अ/ब केंद्रीय सेवा पदांसाठी अर्ज करा. पूर्व, मुख्य आणि मुलाखतीद्वारे निवड.",
			EducationQualification: "Graduate",
			Docs: []string{"Aadhaar Card", "Graduation Degree", "Caste Certificate"},
			MinAge: 21,
			MaxAge: 32,
			LastDate: "2026-07-15",
			Link: "https://upsconline.nic.in",
			GenFee: 200.00,
			ObcFee: 200.00,
			ScStFee: 0.00,
			WomenFee: 0.00,
		},
		{
			Title: "Police Sub-Inspector (PSI) Recruitment 2026",
			TitleHi: "पुलिस सब-इंस्पेक्टर (PSI) भर्ती 2026",
			TitleMr: "पोलीस उपनिरीक्षक (PSI) भरती २०२६",
			Org: "Maharashtra Public Service Commission (MPSC)",
			OrgHi: "महाराष्ट्र लोक सेवा आयोग (MPSC)",
			OrgMr: "महाराष्ट्र लोकसेवा आयोग (MPSC)",
			Desc: "Recruitment of Sub-Inspector in Maharashtra Police Department. Physical standards and physical efficiency test applicable.",
			DescHi: "महाराष्ट्र पुलिस विभाग में सब-इंस्पेक्टर की भर्ती। शारीरिक मानक और शारीरिक दक्षता परीक्षा लागू।",
			DescMr: "महाराष्ट्र पोलीस विभागात उपनिरीक्षक पदाची भरती. शारीरिक पात्रता आणि मैदानी चाचणी लागू.",
			EducationQualification: "Graduate",
			Docs: []string{"Aadhaar Card", "Graduation Degree", "State Domicile Certificate", "Caste Certificate"},
			MinAge: 19,
			MaxAge: 31,
			LastDate: "2026-08-30",
			Link: "https://mpsc.gov.in",
			GenFee: 394.00,
			ObcFee: 294.00,
			ScStFee: 0.00,
			WomenFee: 0.00,
		},
		{
			Title: "Assistant Station Master (ASM)",
			TitleHi: "सहायक स्टेशन मास्टर (ASM)",
			TitleMr: "सहाय्यक स्टेशन मास्टर (ASM)",
			Org: "Railway Recruitment Board (RRB)",
			OrgHi: "रेलवे भर्ती बोर्ड (RRB)",
			OrgMr: "रेल्वे भरती बोर्ड (RRB)",
			Desc: "Excellent career opportunity in Indian Railways for station operations, safety, and train signal coordination.",
			DescHi: "स्टेशन संचालन, सुरक्षा और ट्रेन सिग्नल समन्वय के लिए भारतीय रेलवे में उत्कृष्ट करियर अवसर।",
			DescMr: "स्टेशन ऑपरेशन्स, सुरक्षा आणि ट्रेन सिग्नल समन्वयासाठी भारतीय रेल्वेमध्ये उत्कृष्ट करिअरची संधी.",
			EducationQualification: "Graduate",
			Docs: []string{"Aadhaar Card", "10th Mark Sheet", "Graduation Degree"},
			MinAge: 18,
			MaxAge: 33,
			LastDate: "2026-09-10",
			Link: "https://www.rrcb.gov.in",
			GenFee: 500.00,
			ObcFee: 500.00,
			ScStFee: 250.00,
			WomenFee: 250.00,
		},
		{
			Title: "Technical Assistant & Fireman Grade-A",
			TitleHi: "तकनीकी सहायक और फायरमैन ग्रेड-ए",
			TitleMr: "तांत्रिक सहाय्यक आणि फायरमन ग्रेड-ए",
			Org: "Indian Space Research Organisation (ISRO)",
			OrgHi: "भारतीय अंतरिक्ष अनुसंधान संगठन (ISRO)",
			OrgMr: "भारतीय अंतराळ संशोधन संस्था (ISRO)",
			Desc: "Opportunities for technical diploma and high-school pass holders in specialized research laboratories and launch centers.",
			DescHi: "विशेष अनुसंधान प्रयोगशालाओं और लॉन्च केंद्रों में तकनीकी डिप्लोमा और हाई-स्कूल पास धारकों के लिए अवसर।",
			DescMr: "विशेष संशोधन प्रयोगशाळा आणि प्रक्षेपण केंद्रांमध्ये तांत्रिक डिप्लोमा आणि हायस्कूल उत्तीर्ण धारकांसाठी संधी.",
			EducationQualification: "12th Pass",
			Docs: []string{"Aadhaar Card", "12th Mark Sheet", "Technical Diploma Certificate"},
			MinAge: 18,
			MaxAge: 35,
			LastDate: "2026-06-25",
			Link: "https://www.isro.gov.in",
			GenFee: 250.00,
			ObcFee: 250.00,
			ScStFee: 0.00,
			WomenFee: 0.00,
		},
	}

	for _, j := range jobs {
		query := `
		INSERT INTO government_jobs (
			title, title_hi, title_mr, organization, organization_hi, organization_mr,
			description, description_hi, description_mr, education_qualification, documents_required,
			min_age, max_age, last_date, apply_link, general_fee, obc_fee, sc_st_fee, women_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`

		_, err = DB.Exec(query,
			j.Title, j.TitleHi, j.TitleMr, j.Org, j.OrgHi, j.OrgMr,
			j.Desc, j.DescHi, j.DescMr, j.EducationQualification, pq.StringArray(j.Docs),
			j.MinAge, j.MaxAge, j.LastDate, j.Link, j.GenFee, j.ObcFee, j.ScStFee, j.WomenFee,
		)
		if err != nil {
			return fmt.Errorf("failed to insert job %s: %v", j.Title, err)
		}
	}

	log.Println("Mock government jobs successfully seeded!")
	return nil
}
