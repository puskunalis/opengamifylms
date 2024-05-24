// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: db.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const addUserBadge = `-- name: AddUserBadge :exec
INSERT INTO user_badges (user_id, badge_id)
VALUES ($1, $2)
`

type AddUserBadgeParams struct {
	UserID  int64 `json:"user_id"`
	BadgeID int64 `json:"badge_id"`
}

func (q *Queries) AddUserBadge(ctx context.Context, arg AddUserBadgeParams) error {
	_, err := q.db.Exec(ctx, addUserBadge, arg.UserID, arg.BadgeID)
	return err
}

const checkIfUserAnsweredQuiz = `-- name: CheckIfUserAnsweredQuiz :one
SELECT EXISTS (
    SELECT 1 FROM user_answered_quizzes
    WHERE user_id = $1 AND quiz_id = $2
)
`

type CheckIfUserAnsweredQuizParams struct {
	UserID int64 `json:"user_id"`
	QuizID int64 `json:"quiz_id"`
}

func (q *Queries) CheckIfUserAnsweredQuiz(ctx context.Context, arg CheckIfUserAnsweredQuizParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkIfUserAnsweredQuiz, arg.UserID, arg.QuizID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const checkIfUserCompletedSubmodule = `-- name: CheckIfUserCompletedSubmodule :one
SELECT EXISTS (
    SELECT 1 FROM user_completed_submodules
    WHERE user_id = $1 AND submodule_id = $2
)
`

type CheckIfUserCompletedSubmoduleParams struct {
	UserID      int64 `json:"user_id"`
	SubmoduleID int64 `json:"submodule_id"`
}

func (q *Queries) CheckIfUserCompletedSubmodule(ctx context.Context, arg CheckIfUserCompletedSubmoduleParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkIfUserCompletedSubmodule, arg.UserID, arg.SubmoduleID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createBadge = `-- name: CreateBadge :one
INSERT INTO badges (title, description, icon)
VALUES ($1, $2, $3)
RETURNING id, title, description, icon
`

type CreateBadgeParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (q *Queries) CreateBadge(ctx context.Context, arg CreateBadgeParams) (Badge, error) {
	row := q.db.QueryRow(ctx, createBadge, arg.Title, arg.Description, arg.Icon)
	var i Badge
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Icon,
	)
	return i, err
}

const createCourse = `-- name: CreateCourse :one
INSERT INTO courses (title, description, instructor_id, icon, created_at, updated_at, published, available)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, title, description, instructor_id, icon, created_at, updated_at, published, available
`

type CreateCourseParams struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	InstructorID int64            `json:"instructor_id"`
	Icon         string           `json:"icon"`
	CreatedAt    pgtype.Timestamp `json:"created_at"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
	Published    bool             `json:"published"`
	Available    bool             `json:"available"`
}

func (q *Queries) CreateCourse(ctx context.Context, arg CreateCourseParams) (Course, error) {
	row := q.db.QueryRow(ctx, createCourse,
		arg.Title,
		arg.Description,
		arg.InstructorID,
		arg.Icon,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Published,
		arg.Available,
	)
	var i Course
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.InstructorID,
		&i.Icon,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Published,
		&i.Available,
	)
	return i, err
}

const createElement = `-- name: CreateElement :one
INSERT INTO elements (submodule_id, "type", "content", quiz_id, "order")
VALUES ($1, $2, $3, $4, $5)
RETURNING id, submodule_id, type, content, quiz_id, "order"
`

type CreateElementParams struct {
	SubmoduleID int64       `json:"submodule_id"`
	Type        string      `json:"type"`
	Content     string      `json:"content"`
	QuizID      pgtype.Int8 `json:"quiz_id"`
	Order       int32       `json:"order"`
}

func (q *Queries) CreateElement(ctx context.Context, arg CreateElementParams) (Element, error) {
	row := q.db.QueryRow(ctx, createElement,
		arg.SubmoduleID,
		arg.Type,
		arg.Content,
		arg.QuizID,
		arg.Order,
	)
	var i Element
	err := row.Scan(
		&i.ID,
		&i.SubmoduleID,
		&i.Type,
		&i.Content,
		&i.QuizID,
		&i.Order,
	)
	return i, err
}

const createModule = `-- name: CreateModule :one
INSERT INTO modules (course_id, title, description, "order")
VALUES ($1, $2, $3, $4)
RETURNING id, course_id, title, description, "order"
`

type CreateModuleParams struct {
	CourseID    int64  `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Order       int32  `json:"order"`
}

func (q *Queries) CreateModule(ctx context.Context, arg CreateModuleParams) (Module, error) {
	row := q.db.QueryRow(ctx, createModule,
		arg.CourseID,
		arg.Title,
		arg.Description,
		arg.Order,
	)
	var i Module
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.Title,
		&i.Description,
		&i.Order,
	)
	return i, err
}

const createQuiz = `-- name: CreateQuiz :one
INSERT INTO quizzes (submodule_id, question, question_type, xp_reward)
VALUES ($1, $2, $3, $4)
RETURNING id, submodule_id, question, question_type, xp_reward
`

type CreateQuizParams struct {
	SubmoduleID  int64  `json:"submodule_id"`
	Question     string `json:"question"`
	QuestionType string `json:"question_type"`
	XpReward     int32  `json:"xp_reward"`
}

func (q *Queries) CreateQuiz(ctx context.Context, arg CreateQuizParams) (Quiz, error) {
	row := q.db.QueryRow(ctx, createQuiz,
		arg.SubmoduleID,
		arg.Question,
		arg.QuestionType,
		arg.XpReward,
	)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.SubmoduleID,
		&i.Question,
		&i.QuestionType,
		&i.XpReward,
	)
	return i, err
}

const createQuizAnswer = `-- name: CreateQuizAnswer :one
INSERT INTO quiz_answers (quiz_id, answer_text, is_correct)
VALUES ($1, $2, $3)
RETURNING id, quiz_id, answer_text, is_correct
`

type CreateQuizAnswerParams struct {
	QuizID     int64  `json:"quiz_id"`
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"is_correct"`
}

func (q *Queries) CreateQuizAnswer(ctx context.Context, arg CreateQuizAnswerParams) (QuizAnswer, error) {
	row := q.db.QueryRow(ctx, createQuizAnswer, arg.QuizID, arg.AnswerText, arg.IsCorrect)
	var i QuizAnswer
	err := row.Scan(
		&i.ID,
		&i.QuizID,
		&i.AnswerText,
		&i.IsCorrect,
	)
	return i, err
}

const createSubmodule = `-- name: CreateSubmodule :one
INSERT INTO submodules (module_id, title, xp_reward, "order")
VALUES ($1, $2, $3, $4)
RETURNING id, module_id, title, xp_reward, "order"
`

type CreateSubmoduleParams struct {
	ModuleID int64  `json:"module_id"`
	Title    string `json:"title"`
	XpReward int32  `json:"xp_reward"`
	Order    int32  `json:"order"`
}

func (q *Queries) CreateSubmodule(ctx context.Context, arg CreateSubmoduleParams) (Submodule, error) {
	row := q.db.QueryRow(ctx, createSubmodule,
		arg.ModuleID,
		arg.Title,
		arg.XpReward,
		arg.Order,
	)
	var i Submodule
	err := row.Scan(
		&i.ID,
		&i.ModuleID,
		&i.Title,
		&i.XpReward,
		&i.Order,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (email, full_name, password, is_administrator, is_instructor, xp, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, email, full_name, password, is_administrator, is_instructor, xp, created_at, updated_at
`

type CreateUserParams struct {
	Email           string           `json:"email"`
	FullName        string           `json:"full_name"`
	Password        string           `json:"password"`
	IsAdministrator bool             `json:"is_administrator"`
	IsInstructor    bool             `json:"is_instructor"`
	Xp              int32            `json:"xp"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Email,
		arg.FullName,
		arg.Password,
		arg.IsAdministrator,
		arg.IsInstructor,
		arg.Xp,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.IsAdministrator,
		&i.IsInstructor,
		&i.Xp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createUserCompletedSubmodule = `-- name: CreateUserCompletedSubmodule :exec
INSERT INTO user_completed_submodules (user_id, submodule_id)
VALUES ($1, $2)
`

type CreateUserCompletedSubmoduleParams struct {
	UserID      int64 `json:"user_id"`
	SubmoduleID int64 `json:"submodule_id"`
}

func (q *Queries) CreateUserCompletedSubmodule(ctx context.Context, arg CreateUserCompletedSubmoduleParams) error {
	_, err := q.db.Exec(ctx, createUserCompletedSubmodule, arg.UserID, arg.SubmoduleID)
	return err
}

const deleteBadge = `-- name: DeleteBadge :exec
DELETE FROM badges WHERE id = $1
`

func (q *Queries) DeleteBadge(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteBadge, id)
	return err
}

const deleteCourse = `-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1
`

func (q *Queries) DeleteCourse(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteCourse, id)
	return err
}

const deleteElement = `-- name: DeleteElement :exec
DELETE FROM elements WHERE id = $1
`

func (q *Queries) DeleteElement(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteElement, id)
	return err
}

const deleteEnrollment = `-- name: DeleteEnrollment :exec
DELETE FROM user_courses
WHERE id = $1
`

func (q *Queries) DeleteEnrollment(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteEnrollment, id)
	return err
}

const deleteModule = `-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1
`

func (q *Queries) DeleteModule(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteModule, id)
	return err
}

const deleteQuiz = `-- name: DeleteQuiz :exec
DELETE FROM quizzes WHERE id = $1
`

func (q *Queries) DeleteQuiz(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteQuiz, id)
	return err
}

const deleteQuizAnswer = `-- name: DeleteQuizAnswer :exec
DELETE FROM quiz_answers WHERE id = $1
`

func (q *Queries) DeleteQuizAnswer(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteQuizAnswer, id)
	return err
}

const deleteSubmodule = `-- name: DeleteSubmodule :exec
DELETE FROM submodules WHERE id = $1
`

func (q *Queries) DeleteSubmodule(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteSubmodule, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const enrollUser = `-- name: EnrollUser :exec
INSERT INTO user_courses (user_id, course_id)
VALUES ($1, $2)
`

type EnrollUserParams struct {
	UserID   int64 `json:"user_id"`
	CourseID int64 `json:"course_id"`
}

func (q *Queries) EnrollUser(ctx context.Context, arg EnrollUserParams) error {
	_, err := q.db.Exec(ctx, enrollUser, arg.UserID, arg.CourseID)
	return err
}

const getAllAvailableAndPublishedCourses = `-- name: GetAllAvailableAndPublishedCourses :many
SELECT
    COURSES.ID,
    TITLE,
    DESCRIPTION,
    ICON,
    COURSES.CREATED_AT,
    COURSES.UPDATED_AT,
    USERS.FULL_NAME AS INSTRUCTOR_FULL_NAME,
    COALESCE((SELECT SUM(xp_reward) FROM submodules WHERE module_id IN (SELECT id FROM modules WHERE course_id = COURSES.ID)), 0) AS XP_REWARD
FROM
    COURSES
    JOIN USERS ON USERS.ID = COURSES.INSTRUCTOR_ID
WHERE
    AVAILABLE = TRUE
    AND PUBLISHED = TRUE
`

type GetAllAvailableAndPublishedCoursesRow struct {
	ID                 int64            `json:"id"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	Icon               string           `json:"icon"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
	InstructorFullName string           `json:"instructor_full_name"`
	XpReward           interface{}      `json:"xp_reward"`
}

func (q *Queries) GetAllAvailableAndPublishedCourses(ctx context.Context) ([]GetAllAvailableAndPublishedCoursesRow, error) {
	rows, err := q.db.Query(ctx, getAllAvailableAndPublishedCourses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllAvailableAndPublishedCoursesRow
	for rows.Next() {
		var i GetAllAvailableAndPublishedCoursesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Icon,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.InstructorFullName,
			&i.XpReward,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllBadges = `-- name: GetAllBadges :many
SELECT id, title, description, icon FROM badges
`

func (q *Queries) GetAllBadges(ctx context.Context) ([]Badge, error) {
	rows, err := q.db.Query(ctx, getAllBadges)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Badge
	for rows.Next() {
		var i Badge
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Icon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllCourses = `-- name: GetAllCourses :many
SELECT
    COURSES.ID,
    TITLE,
    DESCRIPTION,
    ICON,
    COURSES.CREATED_AT,
    COURSES.UPDATED_AT,
    USERS.FULL_NAME AS INSTRUCTOR_FULL_NAME,
    COALESCE((SELECT SUM(xp_reward) FROM submodules WHERE module_id IN (SELECT id FROM modules WHERE course_id = COURSES.ID)), 0) AS XP_REWARD
FROM
    COURSES
    JOIN USERS ON USERS.ID = COURSES.INSTRUCTOR_ID
`

type GetAllCoursesRow struct {
	ID                 int64            `json:"id"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	Icon               string           `json:"icon"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
	InstructorFullName string           `json:"instructor_full_name"`
	XpReward           interface{}      `json:"xp_reward"`
}

func (q *Queries) GetAllCourses(ctx context.Context) ([]GetAllCoursesRow, error) {
	rows, err := q.db.Query(ctx, getAllCourses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllCoursesRow
	for rows.Next() {
		var i GetAllCoursesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Icon,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.InstructorFullName,
			&i.XpReward,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllQuizzesBySubmoduleID = `-- name: GetAllQuizzesBySubmoduleID :many
SELECT id, submodule_id, question, question_type, xp_reward FROM quizzes WHERE submodule_id = $1
`

func (q *Queries) GetAllQuizzesBySubmoduleID(ctx context.Context, submoduleID int64) ([]Quiz, error) {
	rows, err := q.db.Query(ctx, getAllQuizzesBySubmoduleID, submoduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Quiz
	for rows.Next() {
		var i Quiz
		if err := rows.Scan(
			&i.ID,
			&i.SubmoduleID,
			&i.Question,
			&i.QuestionType,
			&i.XpReward,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getBadgeByID = `-- name: GetBadgeByID :one
SELECT id, title, description, icon FROM badges WHERE id = $1
`

func (q *Queries) GetBadgeByID(ctx context.Context, id int64) (Badge, error) {
	row := q.db.QueryRow(ctx, getBadgeByID, id)
	var i Badge
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Icon,
	)
	return i, err
}

const getCourseByID = `-- name: GetCourseByID :one
SELECT 
    COURSES.ID,
    TITLE,
    DESCRIPTION,
    ICON,
    COURSES.CREATED_AT,
    COURSES.UPDATED_AT,
    USERS.FULL_NAME AS INSTRUCTOR_FULL_NAME,
    AVAILABLE,
    PUBLISHED,
    USERS.ID AS INSTRUCTOR_ID,
    COALESCE((SELECT SUM(xp_reward) FROM submodules WHERE module_id IN (SELECT id FROM modules WHERE course_id = COURSES.ID)), 0) AS XP_REWARD
FROM
    COURSES
    JOIN USERS ON USERS.ID = COURSES.INSTRUCTOR_ID
WHERE
    COURSES.ID = $1
`

type GetCourseByIDRow struct {
	ID                 int64            `json:"id"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	Icon               string           `json:"icon"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
	InstructorFullName string           `json:"instructor_full_name"`
	Available          bool             `json:"available"`
	Published          bool             `json:"published"`
	InstructorID       int64            `json:"instructor_id"`
	XpReward           interface{}      `json:"xp_reward"`
}

func (q *Queries) GetCourseByID(ctx context.Context, id int64) (GetCourseByIDRow, error) {
	row := q.db.QueryRow(ctx, getCourseByID, id)
	var i GetCourseByIDRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Icon,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.InstructorFullName,
		&i.Available,
		&i.Published,
		&i.InstructorID,
		&i.XpReward,
	)
	return i, err
}

const getCoursesForInstructor = `-- name: GetCoursesForInstructor :many
SELECT id, title, description, instructor_id, icon, created_at, updated_at, published, available FROM courses WHERE instructor_id = $1
`

func (q *Queries) GetCoursesForInstructor(ctx context.Context, instructorID int64) ([]Course, error) {
	rows, err := q.db.Query(ctx, getCoursesForInstructor, instructorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Course
	for rows.Next() {
		var i Course
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.InstructorID,
			&i.Icon,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Published,
			&i.Available,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getElementByID = `-- name: GetElementByID :one
SELECT id, submodule_id, type, content, quiz_id, "order" FROM elements WHERE id = $1
`

func (q *Queries) GetElementByID(ctx context.Context, id int64) (Element, error) {
	row := q.db.QueryRow(ctx, getElementByID, id)
	var i Element
	err := row.Scan(
		&i.ID,
		&i.SubmoduleID,
		&i.Type,
		&i.Content,
		&i.QuizID,
		&i.Order,
	)
	return i, err
}

const getElementsBySubmoduleID = `-- name: GetElementsBySubmoduleID :many
SELECT id, submodule_id, type, content, quiz_id, "order" FROM elements WHERE submodule_id = $1 ORDER BY "order"
`

func (q *Queries) GetElementsBySubmoduleID(ctx context.Context, submoduleID int64) ([]Element, error) {
	rows, err := q.db.Query(ctx, getElementsBySubmoduleID, submoduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Element
	for rows.Next() {
		var i Element
		if err := rows.Scan(
			&i.ID,
			&i.SubmoduleID,
			&i.Type,
			&i.Content,
			&i.QuizID,
			&i.Order,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEnrollment = `-- name: GetEnrollment :one
SELECT id, user_id, course_id
FROM user_courses
WHERE id = $1
`

type GetEnrollmentRow struct {
	ID       int64 `json:"id"`
	UserID   int64 `json:"user_id"`
	CourseID int64 `json:"course_id"`
}

func (q *Queries) GetEnrollment(ctx context.Context, id int64) (GetEnrollmentRow, error) {
	row := q.db.QueryRow(ctx, getEnrollment, id)
	var i GetEnrollmentRow
	err := row.Scan(&i.ID, &i.UserID, &i.CourseID)
	return i, err
}

const getMaxElementOrderForSubmodule = `-- name: GetMaxElementOrderForSubmodule :one
SELECT COALESCE(MAX("order"), -1)::int AS max_order 
FROM elements
WHERE submodule_id = $1
`

func (q *Queries) GetMaxElementOrderForSubmodule(ctx context.Context, submoduleID int64) (int32, error) {
	row := q.db.QueryRow(ctx, getMaxElementOrderForSubmodule, submoduleID)
	var max_order int32
	err := row.Scan(&max_order)
	return max_order, err
}

const getModuleByID = `-- name: GetModuleByID :one
SELECT id, course_id, title, description, "order" FROM modules WHERE id = $1
`

func (q *Queries) GetModuleByID(ctx context.Context, id int64) (Module, error) {
	row := q.db.QueryRow(ctx, getModuleByID, id)
	var i Module
	err := row.Scan(
		&i.ID,
		&i.CourseID,
		&i.Title,
		&i.Description,
		&i.Order,
	)
	return i, err
}

const getModulesByCourseID = `-- name: GetModulesByCourseID :many
SELECT 
    modules.id,
    modules.course_id,
    modules.title,
    modules.description,
    modules."order",
    COALESCE((SELECT SUM(xp_reward) FROM submodules WHERE module_id = modules.id), 0) AS XP_REWARD
FROM modules 
WHERE course_id = $1
ORDER BY "order"
`

type GetModulesByCourseIDRow struct {
	ID          int64       `json:"id"`
	CourseID    int64       `json:"course_id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Order       int32       `json:"order"`
	XpReward    interface{} `json:"xp_reward"`
}

func (q *Queries) GetModulesByCourseID(ctx context.Context, courseID int64) ([]GetModulesByCourseIDRow, error) {
	rows, err := q.db.Query(ctx, getModulesByCourseID, courseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetModulesByCourseIDRow
	for rows.Next() {
		var i GetModulesByCourseIDRow
		if err := rows.Scan(
			&i.ID,
			&i.CourseID,
			&i.Title,
			&i.Description,
			&i.Order,
			&i.XpReward,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuizAnswersByQuizID = `-- name: GetQuizAnswersByQuizID :many
SELECT id, quiz_id, answer_text, is_correct FROM quiz_answers WHERE quiz_id = $1
`

func (q *Queries) GetQuizAnswersByQuizID(ctx context.Context, quizID int64) ([]QuizAnswer, error) {
	rows, err := q.db.Query(ctx, getQuizAnswersByQuizID, quizID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []QuizAnswer
	for rows.Next() {
		var i QuizAnswer
		if err := rows.Scan(
			&i.ID,
			&i.QuizID,
			&i.AnswerText,
			&i.IsCorrect,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getQuizByID = `-- name: GetQuizByID :one
SELECT id, submodule_id, question, question_type, xp_reward FROM quizzes WHERE id = $1
`

func (q *Queries) GetQuizByID(ctx context.Context, id int64) (Quiz, error) {
	row := q.db.QueryRow(ctx, getQuizByID, id)
	var i Quiz
	err := row.Scan(
		&i.ID,
		&i.SubmoduleID,
		&i.Question,
		&i.QuestionType,
		&i.XpReward,
	)
	return i, err
}

const getSubmodule = `-- name: GetSubmodule :one
SELECT id, module_id, title, xp_reward, "order" FROM submodules WHERE id = $1
`

func (q *Queries) GetSubmodule(ctx context.Context, id int64) (Submodule, error) {
	row := q.db.QueryRow(ctx, getSubmodule, id)
	var i Submodule
	err := row.Scan(
		&i.ID,
		&i.ModuleID,
		&i.Title,
		&i.XpReward,
		&i.Order,
	)
	return i, err
}

const getSubmodulesByModuleID = `-- name: GetSubmodulesByModuleID :many
SELECT id, module_id, title, xp_reward, "order" FROM submodules WHERE module_id = $1 ORDER BY "order"
`

func (q *Queries) GetSubmodulesByModuleID(ctx context.Context, moduleID int64) ([]Submodule, error) {
	rows, err := q.db.Query(ctx, getSubmodulesByModuleID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Submodule
	for rows.Next() {
		var i Submodule
		if err := rows.Scan(
			&i.ID,
			&i.ModuleID,
			&i.Title,
			&i.XpReward,
			&i.Order,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTopUsersByXp = `-- name: GetTopUsersByXp :many
SELECT id, full_name, xp
FROM users
ORDER BY xp DESC
LIMIT 10
`

type GetTopUsersByXpRow struct {
	ID       int64  `json:"id"`
	FullName string `json:"full_name"`
	Xp       int32  `json:"xp"`
}

func (q *Queries) GetTopUsersByXp(ctx context.Context) ([]GetTopUsersByXpRow, error) {
	rows, err := q.db.Query(ctx, getTopUsersByXp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTopUsersByXpRow
	for rows.Next() {
		var i GetTopUsersByXpRow
		if err := rows.Scan(&i.ID, &i.FullName, &i.Xp); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserBadges = `-- name: GetUserBadges :many
SELECT b.id, b.title, b.description, b.icon
FROM user_badges ub
JOIN badges b ON ub.badge_id = b.id
WHERE ub.user_id = $1
`

func (q *Queries) GetUserBadges(ctx context.Context, userID int64) ([]Badge, error) {
	rows, err := q.db.Query(ctx, getUserBadges, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Badge
	for rows.Next() {
		var i Badge
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Icon,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, full_name, password, is_administrator, is_instructor, xp, created_at, updated_at
FROM users
WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.Password,
		&i.IsAdministrator,
		&i.IsInstructor,
		&i.Xp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, email, full_name, is_administrator, is_instructor, xp, created_at, updated_at
FROM users
WHERE id = $1
`

type GetUserByIDRow struct {
	ID              int64            `json:"id"`
	Email           string           `json:"email"`
	FullName        string           `json:"full_name"`
	IsAdministrator bool             `json:"is_administrator"`
	IsInstructor    bool             `json:"is_instructor"`
	Xp              int32            `json:"xp"`
	CreatedAt       pgtype.Timestamp `json:"created_at"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) GetUserByID(ctx context.Context, id int64) (GetUserByIDRow, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i GetUserByIDRow
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.IsAdministrator,
		&i.IsInstructor,
		&i.Xp,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserCompletedSubmodulesByUserID = `-- name: GetUserCompletedSubmodulesByUserID :many
SELECT id, user_id, submodule_id, completed_at FROM user_completed_submodules
WHERE user_id = $1
`

func (q *Queries) GetUserCompletedSubmodulesByUserID(ctx context.Context, userID int64) ([]UserCompletedSubmodule, error) {
	rows, err := q.db.Query(ctx, getUserCompletedSubmodulesByUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UserCompletedSubmodule
	for rows.Next() {
		var i UserCompletedSubmodule
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SubmoduleID,
			&i.CompletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserCourses = `-- name: GetUserCourses :many
SELECT 
    courses.id AS id,
    title, 
    description, 
    icon, 
    courses.created_at, 
    courses.updated_at, 
    published, 
    available, 
    users.full_name AS instructor_full_name,
    COALESCE((SELECT SUM(xp_reward) FROM submodules WHERE module_id IN (SELECT id FROM modules WHERE course_id = courses.id)), 0) AS XP_REWARD
FROM user_courses
JOIN courses on user_courses.course_id = courses.id
JOIN users on users.id = courses.instructor_id
WHERE user_id = $1
`

type GetUserCoursesRow struct {
	ID                 int64            `json:"id"`
	Title              string           `json:"title"`
	Description        string           `json:"description"`
	Icon               string           `json:"icon"`
	CreatedAt          pgtype.Timestamp `json:"created_at"`
	UpdatedAt          pgtype.Timestamp `json:"updated_at"`
	Published          bool             `json:"published"`
	Available          bool             `json:"available"`
	InstructorFullName string           `json:"instructor_full_name"`
	XpReward           interface{}      `json:"xp_reward"`
}

func (q *Queries) GetUserCourses(ctx context.Context, userID int64) ([]GetUserCoursesRow, error) {
	rows, err := q.db.Query(ctx, getUserCourses, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserCoursesRow
	for rows.Next() {
		var i GetUserCoursesRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Icon,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Published,
			&i.Available,
			&i.InstructorFullName,
			&i.XpReward,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const setQuizAnsweredByUser = `-- name: SetQuizAnsweredByUser :exec
INSERT INTO user_answered_quizzes (user_id, quiz_id)
VALUES ($1, $2)
`

type SetQuizAnsweredByUserParams struct {
	UserID int64 `json:"user_id"`
	QuizID int64 `json:"quiz_id"`
}

func (q *Queries) SetQuizAnsweredByUser(ctx context.Context, arg SetQuizAnsweredByUserParams) error {
	_, err := q.db.Exec(ctx, setQuizAnsweredByUser, arg.UserID, arg.QuizID)
	return err
}

const setUserCourseComplete = `-- name: SetUserCourseComplete :exec
UPDATE user_courses
SET completed = true
WHERE course_id = $1
`

func (q *Queries) SetUserCourseComplete(ctx context.Context, courseID int64) error {
	_, err := q.db.Exec(ctx, setUserCourseComplete, courseID)
	return err
}

const updateBadge = `-- name: UpdateBadge :exec
UPDATE badges
SET title = $1, icon = $2
WHERE id = $3
`

type UpdateBadgeParams struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	ID    int64  `json:"id"`
}

func (q *Queries) UpdateBadge(ctx context.Context, arg UpdateBadgeParams) error {
	_, err := q.db.Exec(ctx, updateBadge, arg.Title, arg.Icon, arg.ID)
	return err
}

const updateCourse = `-- name: UpdateCourse :one
UPDATE courses
SET title = $1, description = $2, instructor_id = $3, updated_at = $4, published = $5, available = $6
WHERE id = $7
RETURNING id, title, description, instructor_id, icon, created_at, updated_at, published, available
`

type UpdateCourseParams struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	InstructorID int64            `json:"instructor_id"`
	UpdatedAt    pgtype.Timestamp `json:"updated_at"`
	Published    bool             `json:"published"`
	Available    bool             `json:"available"`
	ID           int64            `json:"id"`
}

func (q *Queries) UpdateCourse(ctx context.Context, arg UpdateCourseParams) (Course, error) {
	row := q.db.QueryRow(ctx, updateCourse,
		arg.Title,
		arg.Description,
		arg.InstructorID,
		arg.UpdatedAt,
		arg.Published,
		arg.Available,
		arg.ID,
	)
	var i Course
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.InstructorID,
		&i.Icon,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Published,
		&i.Available,
	)
	return i, err
}

const updateElementOrderBatch = `-- name: UpdateElementOrderBatch :exec
UPDATE elements
SET "order" = new_orders.order
FROM (
    SELECT unnest($2::int[]) AS id, unnest($3::int[]) AS order
) AS new_orders
WHERE elements.id = new_orders.id AND elements.submodule_id = $1
`

type UpdateElementOrderBatchParams struct {
	SubmoduleID int64   `json:"submodule_id"`
	Ids         []int32 `json:"ids"`
	Orders      []int32 `json:"orders"`
}

func (q *Queries) UpdateElementOrderBatch(ctx context.Context, arg UpdateElementOrderBatchParams) error {
	_, err := q.db.Exec(ctx, updateElementOrderBatch, arg.SubmoduleID, arg.Ids, arg.Orders)
	return err
}

const updateModuleOrderBatch = `-- name: UpdateModuleOrderBatch :exec
UPDATE modules
SET "order" = new_orders.order
FROM (
    SELECT unnest($2::int[]) AS id, unnest($3::int[]) AS order
) AS new_orders
WHERE modules.id = new_orders.id AND modules.course_id = $1
`

type UpdateModuleOrderBatchParams struct {
	CourseID int64   `json:"course_id"`
	Ids      []int32 `json:"ids"`
	Orders   []int32 `json:"orders"`
}

func (q *Queries) UpdateModuleOrderBatch(ctx context.Context, arg UpdateModuleOrderBatchParams) error {
	_, err := q.db.Exec(ctx, updateModuleOrderBatch, arg.CourseID, arg.Ids, arg.Orders)
	return err
}

const updateQuiz = `-- name: UpdateQuiz :exec
UPDATE quizzes
SET question = $1, question_type = $2
WHERE id = $3
`

type UpdateQuizParams struct {
	Question     string `json:"question"`
	QuestionType string `json:"question_type"`
	ID           int64  `json:"id"`
}

func (q *Queries) UpdateQuiz(ctx context.Context, arg UpdateQuizParams) error {
	_, err := q.db.Exec(ctx, updateQuiz, arg.Question, arg.QuestionType, arg.ID)
	return err
}

const updateQuizAnswer = `-- name: UpdateQuizAnswer :exec
UPDATE quiz_answers
SET answer_text = $1, is_correct = $2
WHERE id = $3
`

type UpdateQuizAnswerParams struct {
	AnswerText string `json:"answer_text"`
	IsCorrect  bool   `json:"is_correct"`
	ID         int64  `json:"id"`
}

func (q *Queries) UpdateQuizAnswer(ctx context.Context, arg UpdateQuizAnswerParams) error {
	_, err := q.db.Exec(ctx, updateQuizAnswer, arg.AnswerText, arg.IsCorrect, arg.ID)
	return err
}

const updateSubmoduleOrderBatch = `-- name: UpdateSubmoduleOrderBatch :exec
UPDATE submodules
SET "order" = new_orders.order
FROM (
    SELECT unnest($2::int[]) AS id, unnest($3::int[]) AS order
) AS new_orders
WHERE submodules.id = new_orders.id AND submodules.module_id = $1
`

type UpdateSubmoduleOrderBatchParams struct {
	ModuleID int64   `json:"module_id"`
	Ids      []int32 `json:"ids"`
	Orders   []int32 `json:"orders"`
}

func (q *Queries) UpdateSubmoduleOrderBatch(ctx context.Context, arg UpdateSubmoduleOrderBatchParams) error {
	_, err := q.db.Exec(ctx, updateSubmoduleOrderBatch, arg.ModuleID, arg.Ids, arg.Orders)
	return err
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET email = $1, full_name = $2, password = $3, is_administrator = $4, is_instructor = $5, xp = $6, updated_at = $7
WHERE id = $8
`

type UpdateUserParams struct {
	Email           string           `json:"email"`
	FullName        string           `json:"full_name"`
	Password        string           `json:"password"`
	IsAdministrator bool             `json:"is_administrator"`
	IsInstructor    bool             `json:"is_instructor"`
	Xp              int32            `json:"xp"`
	UpdatedAt       pgtype.Timestamp `json:"updated_at"`
	ID              int64            `json:"id"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.Exec(ctx, updateUser,
		arg.Email,
		arg.FullName,
		arg.Password,
		arg.IsAdministrator,
		arg.IsInstructor,
		arg.Xp,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}

const updateUserXp = `-- name: UpdateUserXp :exec
UPDATE users
SET xp = $1
WHERE id = $2
`

type UpdateUserXpParams struct {
	Xp int32 `json:"xp"`
	ID int64 `json:"id"`
}

func (q *Queries) UpdateUserXp(ctx context.Context, arg UpdateUserXpParams) error {
	_, err := q.db.Exec(ctx, updateUserXp, arg.Xp, arg.ID)
	return err
}
