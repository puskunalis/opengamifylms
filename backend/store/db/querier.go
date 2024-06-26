// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"context"
)

type Querier interface {
	AddUserBadge(ctx context.Context, arg AddUserBadgeParams) error
	CheckIfUserAnsweredQuiz(ctx context.Context, arg CheckIfUserAnsweredQuizParams) (bool, error)
	CheckIfUserCompletedSubmodule(ctx context.Context, arg CheckIfUserCompletedSubmoduleParams) (bool, error)
	CreateBadge(ctx context.Context, arg CreateBadgeParams) (Badge, error)
	CreateCourse(ctx context.Context, arg CreateCourseParams) (Course, error)
	CreateElement(ctx context.Context, arg CreateElementParams) (Element, error)
	CreateModule(ctx context.Context, arg CreateModuleParams) (Module, error)
	CreateQuiz(ctx context.Context, arg CreateQuizParams) (Quiz, error)
	CreateQuizAnswer(ctx context.Context, arg CreateQuizAnswerParams) (QuizAnswer, error)
	CreateSubmodule(ctx context.Context, arg CreateSubmoduleParams) (Submodule, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	CreateUserCompletedSubmodule(ctx context.Context, arg CreateUserCompletedSubmoduleParams) error
	DeleteBadge(ctx context.Context, id int64) error
	DeleteCourse(ctx context.Context, id int64) error
	DeleteElement(ctx context.Context, id int64) error
	DeleteEnrollment(ctx context.Context, id int64) error
	DeleteModule(ctx context.Context, id int64) error
	DeleteQuiz(ctx context.Context, id int64) error
	DeleteQuizAnswer(ctx context.Context, id int64) error
	DeleteSubmodule(ctx context.Context, id int64) error
	DeleteUser(ctx context.Context, id int64) error
	EnrollUser(ctx context.Context, arg EnrollUserParams) error
	GetAllAvailableAndPublishedCourses(ctx context.Context) ([]GetAllAvailableAndPublishedCoursesRow, error)
	GetAllBadges(ctx context.Context) ([]Badge, error)
	GetAllCourses(ctx context.Context) ([]GetAllCoursesRow, error)
	GetAllQuizzesBySubmoduleID(ctx context.Context, submoduleID int64) ([]Quiz, error)
	GetBadgeByID(ctx context.Context, id int64) (Badge, error)
	GetCourseByID(ctx context.Context, id int64) (GetCourseByIDRow, error)
	GetCoursesForInstructor(ctx context.Context, instructorID int64) ([]Course, error)
	GetElementByID(ctx context.Context, id int64) (Element, error)
	GetElementsBySubmoduleID(ctx context.Context, submoduleID int64) ([]Element, error)
	GetEnrollment(ctx context.Context, id int64) (GetEnrollmentRow, error)
	GetMaxElementOrderForSubmodule(ctx context.Context, submoduleID int64) (int32, error)
	GetModuleByID(ctx context.Context, id int64) (Module, error)
	GetModulesByCourseID(ctx context.Context, courseID int64) ([]GetModulesByCourseIDRow, error)
	GetQuizAnswersByQuizID(ctx context.Context, quizID int64) ([]QuizAnswer, error)
	GetQuizByID(ctx context.Context, id int64) (Quiz, error)
	GetSubmodule(ctx context.Context, id int64) (Submodule, error)
	GetSubmodulesByModuleID(ctx context.Context, moduleID int64) ([]Submodule, error)
	GetTopUsersByXp(ctx context.Context) ([]GetTopUsersByXpRow, error)
	GetUserBadges(ctx context.Context, userID int64) ([]Badge, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error)
	GetUserCompletedSubmodulesByUserID(ctx context.Context, userID int64) ([]UserCompletedSubmodule, error)
	GetUserCourses(ctx context.Context, userID int64) ([]GetUserCoursesRow, error)
	SetQuizAnsweredByUser(ctx context.Context, arg SetQuizAnsweredByUserParams) error
	SetUserCourseComplete(ctx context.Context, courseID int64) error
	UpdateBadge(ctx context.Context, arg UpdateBadgeParams) error
	UpdateCourse(ctx context.Context, arg UpdateCourseParams) (Course, error)
	UpdateElementOrderBatch(ctx context.Context, arg UpdateElementOrderBatchParams) error
	UpdateModuleOrderBatch(ctx context.Context, arg UpdateModuleOrderBatchParams) error
	UpdateQuiz(ctx context.Context, arg UpdateQuizParams) error
	UpdateQuizAnswer(ctx context.Context, arg UpdateQuizAnswerParams) error
	UpdateSubmoduleOrderBatch(ctx context.Context, arg UpdateSubmoduleOrderBatchParams) error
	UpdateUser(ctx context.Context, arg UpdateUserParams) error
	UpdateUserXp(ctx context.Context, arg UpdateUserXpParams) error
}

var _ Querier = (*Queries)(nil)
