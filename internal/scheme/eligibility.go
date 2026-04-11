package scheme

import (
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job"
	"github.com/ONEST-Network/Job-Manager-Adapter/pkg/database/mongodb/job-application"
)

// ValidateEligibility compares applicant details against scheme requirements
// Returns an eligibility score (0.0 to 1.0) and whether they are eligible
func ValidateEligibility(criteria job.Eligibility, applicant jobapplication.ApplicantDetails) (float64, bool) {
	score := 0.0
	totalWeight := 3.0

	// 1. Gender check
	if criteria.Gender == job.GenderAny || string(criteria.Gender) == applicant.Gender {
		score += 1.0
	}

	// 2. Experience check
	if applicant.Experience.Years >= criteria.YearsOfExperience {
		score += 1.0
	} else if criteria.YearsOfExperience > 0 {
		score += float64(applicant.Experience.Years) / float64(criteria.YearsOfExperience)
	}

	// 3. Academic Qualification check
	if isQualified(criteria.AcademicQualification, applicant) {
		score += 1.0
	}

	finalScore := score / totalWeight
	isEligible := finalScore >= 0.7 // Threshold for eligibility

	return finalScore, isEligible
}

func isQualified(required job.AcademicQualification, applicant jobapplication.ApplicantDetails) bool {
	// Simplified qualification hierarchy check
	qualMap := map[job.AcademicQualification]int{
		job.AcademicQualificationNone:         0,
		job.AcademicQualificationClassX:       1,
		job.AcademicQualificationClassXII:     2,
		job.AcademicQualificationDiploma:      3,
		job.AcademicQualificationGraduate:     4,
		job.AcademicQualificationPostGraduate: 5,
	}

	// We don't have applicant's qualification directly in ApplicantDetails,
	// but we could infer it or assume it's part of a complete profile.
	// For now, let's assume applicant has at least some qualification if age > 18
	applicantQual := job.AcademicQualificationNone
	if applicant.Age >= 22 {
		applicantQual = job.AcademicQualificationGraduate
	} else if applicant.Age >= 18 {
		applicantQual = job.AcademicQualificationClassXII
	}

	return qualMap[applicantQual] >= qualMap[required]
}
