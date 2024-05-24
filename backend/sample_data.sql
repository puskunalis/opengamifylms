-- Insert Instructor: John Doe
INSERT INTO users (email, full_name, password, is_instructor)
VALUES ('john@example.com', 'John Doe', '$2a$10$62fZEu7Un1exKq4vSuXykeaAyqQL2GbTlNcIIqFULmqijLVPz5ugS', true);

-- Get the inserted user ID for John Doe
SELECT currval(pg_get_serial_sequence('users', 'id')) AS john_doe_user_id;

-- Insert Instructor: Jane Smith
INSERT INTO users (email, full_name, password, is_instructor)
VALUES ('jane@example.com', 'Jane Smith', '$2a$10$62fZEu7Un1exKq4vSuXykeaAyqQL2GbTlNcIIqFULmqijLVPz5ugS', true);

-- Get the inserted user ID for Jane Smith
SELECT currval(pg_get_serial_sequence('users', 'id')) AS jane_smith_user_id;

-- Insert Course: Introduction to Programming
INSERT INTO courses (title, description, instructor_id, icon, published, available)
VALUES ('Introduction to Programming', 'Learn the basics of programming concepts and algorithms.', 1, 'https://picsum.photos/seed/c1/768/432', true, true);

-- Get the inserted course ID
SELECT currval(pg_get_serial_sequence('courses', 'id')) AS intro_to_programming_course_id;

-- Insert Module: HTML
INSERT INTO modules (course_id, title, description, "order")
VALUES ((SELECT currval(pg_get_serial_sequence('courses', 'id'))), 'HTML', 'A dive into the simplicity of HTML.', 1);

-- Get the inserted module ID
SELECT currval(pg_get_serial_sequence('modules', 'id')) AS html_module_id;

-- Insert Submodule: Understanding HTML Syntax
INSERT INTO submodules (module_id, title, xp_reward, "order")
VALUES ((SELECT currval(pg_get_serial_sequence('modules', 'id'))), 'Understanding HTML Syntax', 20, 1);

-- Get the inserted submodule ID
SELECT currval(pg_get_serial_sequence('submodules', 'id')) AS understanding_html_syntax_submodule_id;

-- Insert Elements for Understanding HTML Syntax submodule
INSERT INTO elements (submodule_id, type, content, "order", quiz_id)
VALUES
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<h1>Syntax</h1>', 1, NULL),
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<p>HTML has cool syntax.</p>', 2, NULL),
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<p>HTML has great syntax, actually.</p>', 3, NULL);
--    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'quiz_single_choice', '{"question":"What is the capital of France?","answers":["London","Paris","Berlin","Madrid"],"correct_answer":"Paris"}', 4),
--    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'quiz_single_choice', '{"question":"What is the capital of Estonia?","answers":["Riga","Tallinn","Helsinki"],"correct_answer":"Tallinn"}', 5),
--    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'quiz_multiple_choice', '{"question":"Favourite food of Lithuanians?","answers":["Cepelinai","Niekas","Saltibarsciai"],"correct_answers":["Cepelinai","Saltibarsciai"]}', 6),
--    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'video', 'http://127.0.0.1:9000/static/20240127143026_212663.MP4', 7);

-- Insert Submodule: HTML inner working
INSERT INTO submodules (module_id, title, xp_reward, "order")
VALUES ((SELECT currval(pg_get_serial_sequence('modules', 'id'))), 'HTML inner working', 15, 2);

-- Get the inserted submodule ID
SELECT currval(pg_get_serial_sequence('submodules', 'id')) AS html_inner_working_submodule_id;

-- Insert Elements for HTML inner working submodule
INSERT INTO elements (submodule_id, type, content, "order", quiz_id)
VALUES
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<h1>HTML inside</h1>', 1, NULL),
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<h2>Nobody knows how it really works, actually.</h2>', 2, NULL);
--    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'quiz_single_choice', '{"question":"What is the capital of Lithuania?","answers":["Vilnius","Paris","Madrid"],"correct_answer":"Vilnius"}', 3),
--    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'quiz_single_choice', '{"question":"What is the capital of Latvia?","answers":["Riga","Vilnius","Tallinn"],"correct_answer":"Riga"}', 4);

-- Insert Module: CSS
INSERT INTO modules (course_id, title, description, "order")
VALUES ((SELECT currval(pg_get_serial_sequence('courses', 'id'))), 'CSS', 'Learn the difficult basics of CSS.', 2);

-- Get the inserted module ID
SELECT currval(pg_get_serial_sequence('modules', 'id')) AS css_module_id;

-- Insert Submodule: Getting CSS
INSERT INTO submodules (module_id, title, xp_reward, "order")
VALUES ((SELECT currval(pg_get_serial_sequence('modules', 'id'))), 'Getting CSS', 5, 1);

-- Get the inserted submodule ID
SELECT currval(pg_get_serial_sequence('submodules', 'id')) AS getting_css_submodule_id;

-- Insert Elements for Getting CSS submodule
INSERT INTO elements (submodule_id, type, content, "order")
VALUES
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<h1>CSS</h1>', 1),
    ((SELECT currval(pg_get_serial_sequence('submodules', 'id'))), 'html', '<p>CSS is difficult.</p>', 2);

-- Insert additional courses
INSERT INTO courses (title, description, instructor_id, icon, published, available)
VALUES ('Web Development Fundamentals', 'Explore the fundamentals of web development using HTML, CSS, and JavaScript.', 2, 'https://picsum.photos/seed/c2/768/432', true, true);
INSERT INTO courses (title, description, instructor_id, icon, published, available)
VALUES ('Data Structures and Algorithms', 'Master the essential data structures and algorithms for efficient problem-solving.', 1, 'https://picsum.photos/seed/c3/768/432', true, true);

-- Insert badges
INSERT INTO badges (title, description, icon)
VALUES
    ('Big Learner', 'You go above and beyond!', 'https://picsum.photos/seed/b1/100'),
    ('Persistent Learner', 'You are persistent!', 'https://picsum.photos/seed/b2/100');

-- Insert a user badge
INSERT INTO user_badges (user_id, badge_id)
VALUES ((SELECT id FROM users WHERE full_name = 'John Doe'), (SELECT id FROM badges WHERE title = 'Big Learner'));
