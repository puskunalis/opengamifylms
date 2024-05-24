import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, TextField, Button } from '@mui/material';
import Layout from './Layout';
import InstructorCourseList from './InstructorCourseList';

const CourseBuilder = ({ user }) => {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [instructorCourses, setInstructorCourses] = useState([]);

    useEffect(() => {
        const fetchInstructorCourses = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/courses`, {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    },
                });

                if (response.ok) {
                    const data = await response.json();
                    setInstructorCourses(data);
                } else {
                    console.error('Failed to fetch instructor courses');
                }
            } catch (error) {
                console.error('Error fetching instructor courses:', error);
            }
        };

        fetchInstructorCourses();
    }, [user.id]);

    const handleTitleChange = (event) => {
        setTitle(event.target.value);
    };

    const handleDescriptionChange = (event) => {
        setDescription(event.target.value);
    };

    const handleSaveCourse = async () => {
        if (title.trim() !== '' && description.trim() !== '') {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    },
                    body: JSON.stringify({
                        title: title,
                        description: description,
                    }),
                });

                if (response.ok) {
                    const newCourse = await response.json();
                    if (instructorCourses) {
                        setInstructorCourses([...instructorCourses, newCourse]);
                    } else {
                        setInstructorCourses([newCourse]);
                    }
                    setTitle('');
                    setDescription('');
                    alert('Course saved successfully!');
                } else {
                    console.error('Failed to save course');
                }
            } catch (error) {
                console.error('Error saving course:', error);
            }
        }
    };

    return (
        <Layout user={user}>
            <Container maxWidth="lg" className="my-8">
                <Box mb={4}>
                    <Typography variant="h4" component="h1" gutterBottom>
                        Course Builder
                    </Typography>

                    <TextField
                        label="Course Title"
                        value={title}
                        onChange={handleTitleChange}
                        fullWidth
                        margin="normal"
                    />

                    <TextField
                        label="Course Description"
                        value={description}
                        onChange={handleDescriptionChange}
                        fullWidth
                        multiline
                        rows={4}
                        margin="normal"
                    />

                    <Button variant="contained" color="primary" onClick={handleSaveCourse}>
                        Save Course
                    </Button>
                </Box>
                {instructorCourses && <Box mt={4}>
                    <Typography variant="h5" component="h2" gutterBottom>
                        Your Courses
                    </Typography>
                    <InstructorCourseList courses={instructorCourses} setCourses={setInstructorCourses} />
                </Box>}
            </Container>
        </Layout>
    );
};

export default CourseBuilder;