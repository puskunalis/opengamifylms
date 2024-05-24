CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    full_name TEXT NOT NULL,
    password TEXT NOT NULL,
    is_administrator BOOLEAN NOT NULL DEFAULT false,
    is_instructor BOOLEAN NOT NULL DEFAULT false,
    xp INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    instructor_id BIGSERIAL REFERENCES users(id) NOT NULL,
    icon TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published BOOLEAN NOT NULL DEFAULT false,
    available BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE user_courses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users(id) NOT NULL,
    course_id BIGSERIAL REFERENCES courses(id) NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (user_id, course_id)
);

CREATE TABLE modules (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGSERIAL REFERENCES courses(id) NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    "order" INTEGER NOT NULL
);

CREATE TABLE submodules (
    id BIGSERIAL PRIMARY KEY,
    module_id BIGSERIAL REFERENCES modules(id) NOT NULL,
    title TEXT NOT NULL,
    xp_reward INTEGER NOT NULL,
    "order" INTEGER NOT NULL
);

CREATE TABLE quizzes (
    id BIGSERIAL PRIMARY KEY,
    submodule_id BIGSERIAL REFERENCES submodules(id) NOT NULL,
    question TEXT NOT NULL,
    question_type TEXT NOT NULL CHECK (question_type IN ('single', 'multiple')),
    xp_reward INTEGER NOT NULL
);

CREATE TABLE elements (
    id BIGSERIAL PRIMARY KEY,
    submodule_id BIGSERIAL REFERENCES submodules(id) NOT NULL,
    type TEXT NOT NULL,
    content TEXT NOT NULL,
    quiz_id BIGINT REFERENCES quizzes(id),
    "order" INTEGER NOT NULL
);

CREATE TABLE badges (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    icon TEXT NOT NULL
);

CREATE TABLE user_badges (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users(id) NOT NULL,
    badge_id BIGSERIAL REFERENCES badges(id) NOT NULL,
    UNIQUE (user_id, badge_id)
);

CREATE TABLE quiz_answers (
    id BIGSERIAL PRIMARY KEY,
    quiz_id BIGSERIAL REFERENCES quizzes(id) NOT NULL,
    answer_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL
);

CREATE TABLE user_answered_quizzes (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users(id) NOT NULL,
    quiz_id BIGSERIAL REFERENCES quizzes(id) NOT NULL,
    UNIQUE (user_id, quiz_id)
);

CREATE TABLE user_completed_submodules (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGSERIAL REFERENCES users(id) NOT NULL,
    submodule_id BIGSERIAL REFERENCES submodules(id) NOT NULL,
    completed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, submodule_id)
);
