-- name: CreateCourse :one
INSERT INTO courses (title, description, instructor_id, icon, created_at, updated_at, published, available)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetCourseByID :one
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
    COURSES.ID = $1;

-- name: GetAllCourses :many
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
    JOIN USERS ON USERS.ID = COURSES.INSTRUCTOR_ID;

-- name: GetAllAvailableAndPublishedCourses :many
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
    AND PUBLISHED = TRUE;

-- name: UpdateCourse :one
UPDATE courses
SET title = $1, description = $2, instructor_id = $3, updated_at = $4, published = $5, available = $6
WHERE id = $7
RETURNING *;

-- name: GetCoursesForInstructor :many
SELECT * FROM courses WHERE instructor_id = $1;

-- name: DeleteCourse :exec
DELETE FROM courses WHERE id = $1;

-- name: GetModulesByCourseID :many
SELECT 
    modules.id,
    modules.course_id,
    modules.title,
    modules.description,
    modules."order",
    COALESCE((SELECT SUM(xp_reward) FROM submodules WHERE module_id = modules.id), 0) AS XP_REWARD
FROM modules 
WHERE course_id = $1
ORDER BY "order";

-- name: CreateModule :one
INSERT INTO modules (course_id, title, description, "order")
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteModule :exec
DELETE FROM modules WHERE id = $1;

-- name: UpdateModuleOrderBatch :exec
UPDATE modules
SET "order" = new_orders.order
FROM (
    SELECT unnest(@ids::int[]) AS id, unnest(@orders::int[]) AS order
) AS new_orders
WHERE modules.id = new_orders.id AND modules.course_id = @course_id;

-- name: CreateSubmodule :one
INSERT INTO submodules (module_id, title, xp_reward, "order")
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: DeleteSubmodule :exec
DELETE FROM submodules WHERE id = $1;

-- name: UpdateSubmoduleOrderBatch :exec
UPDATE submodules
SET "order" = new_orders.order
FROM (
    SELECT unnest(@ids::int[]) AS id, unnest(@orders::int[]) AS order
) AS new_orders
WHERE submodules.id = new_orders.id AND submodules.module_id = @module_id;

-- name: CreateElement :one
INSERT INTO elements (submodule_id, "type", "content", quiz_id, "order")
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetElementByID :one
SELECT * FROM elements WHERE id = $1;

-- name: GetMaxElementOrderForSubmodule :one
SELECT COALESCE(MAX("order"), -1)::int AS max_order 
FROM elements
WHERE submodule_id = $1;

-- name: GetElementsBySubmoduleID :many
SELECT * FROM elements WHERE submodule_id = $1 ORDER BY "order";

-- name: DeleteElement :exec
DELETE FROM elements WHERE id = $1;

-- name: UpdateElementOrderBatch :exec
UPDATE elements
SET "order" = new_orders.order
FROM (
    SELECT unnest(@ids::int[]) AS id, unnest(@orders::int[]) AS order
) AS new_orders
WHERE elements.id = new_orders.id AND elements.submodule_id = @submodule_id;

-- name: GetModuleByID :one
SELECT * FROM modules WHERE id = $1;

-- name: GetSubmodule :one
SELECT * FROM submodules WHERE id = $1;

-- name: GetSubmodulesByModuleID :many
SELECT * FROM submodules WHERE module_id = $1 ORDER BY "order";

-- name: CreateUser :one
INSERT INTO users (email, full_name, password, is_administrator, is_instructor, xp, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByID :one
SELECT id, email, full_name, is_administrator, is_instructor, xp, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserCourses :many
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
WHERE user_id = $1;

-- name: GetUserByEmail :one
SELECT id, email, full_name, password, is_administrator, is_instructor, xp, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UpdateUser :exec
UPDATE users
SET email = $1, full_name = $2, password = $3, is_administrator = $4, is_instructor = $5, xp = $6, updated_at = $7
WHERE id = $8;

-- name: UpdateUserXp :exec
UPDATE users
SET xp = $1
WHERE id = $2;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: EnrollUser :exec
INSERT INTO user_courses (user_id, course_id)
VALUES ($1, $2);

-- name: GetTopUsersByXp :many
SELECT id, full_name, xp
FROM users
ORDER BY xp DESC
LIMIT 10;

-- name: GetEnrollment :one
SELECT id, user_id, course_id
FROM user_courses
WHERE id = $1;

-- name: DeleteEnrollment :exec
DELETE FROM user_courses
WHERE id = $1;

-- name: SetUserCourseComplete :exec
UPDATE user_courses
SET completed = true
WHERE course_id = $1;

-- name: CreateQuiz :one
INSERT INTO quizzes (submodule_id, question, question_type, xp_reward)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetQuizByID :one
SELECT * FROM quizzes WHERE id = $1;

-- name: GetAllQuizzesBySubmoduleID :many
SELECT * FROM quizzes WHERE submodule_id = $1;

-- name: UpdateQuiz :exec
UPDATE quizzes
SET question = $1, question_type = $2
WHERE id = $3;

-- name: DeleteQuiz :exec
DELETE FROM quizzes WHERE id = $1;

-- name: CreateQuizAnswer :one
INSERT INTO quiz_answers (quiz_id, answer_text, is_correct)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetQuizAnswersByQuizID :many
SELECT * FROM quiz_answers WHERE quiz_id = $1;

-- name: UpdateQuizAnswer :exec
UPDATE quiz_answers
SET answer_text = $1, is_correct = $2
WHERE id = $3;

-- name: DeleteQuizAnswer :exec
DELETE FROM quiz_answers WHERE id = $1;

-- name: SetQuizAnsweredByUser :exec
INSERT INTO user_answered_quizzes (user_id, quiz_id)
VALUES ($1, $2);

-- name: CheckIfUserAnsweredQuiz :one
SELECT EXISTS (
    SELECT 1 FROM user_answered_quizzes
    WHERE user_id = $1 AND quiz_id = $2
);

-- name: CreateBadge :one
INSERT INTO badges (title, description, icon)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBadgeByID :one
SELECT * FROM badges WHERE id = $1;

-- name: GetAllBadges :many
SELECT * FROM badges;

-- name: UpdateBadge :exec
UPDATE badges
SET title = $1, icon = $2
WHERE id = $3;

-- name: DeleteBadge :exec
DELETE FROM badges WHERE id = $1;

-- name: AddUserBadge :exec
INSERT INTO user_badges (user_id, badge_id)
VALUES ($1, $2);

-- name: GetUserBadges :many
SELECT b.id, b.title, b.description, b.icon
FROM user_badges ub
JOIN badges b ON ub.badge_id = b.id
WHERE ub.user_id = $1;

-- name: CreateUserCompletedSubmodule :exec
INSERT INTO user_completed_submodules (user_id, submodule_id)
VALUES ($1, $2);

-- name: CheckIfUserCompletedSubmodule :one
SELECT EXISTS (
    SELECT 1 FROM user_completed_submodules
    WHERE user_id = $1 AND submodule_id = $2
);

-- name: GetUserCompletedSubmodulesByUserID :many
SELECT * FROM user_completed_submodules
WHERE user_id = $1;