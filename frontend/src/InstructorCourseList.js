import React from 'react';
import { List, Grid, Card, CardActionArea, Typography, CardContent, CardMedia, Chip, CardActions, Button } from '@mui/material';
import { Link } from 'react-router-dom';

const InstructorCourseList = ({ courses, setCourses }) => {
    const [customSystemSettings, setCustomSystemSettings] = React.useState(null);

    React.useEffect(() => {
        const fetchCustomSystemSettings = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/customSystemSettings`);
                const customSystemSettingsData = await response.json();
                setCustomSystemSettings(customSystemSettingsData);
            } catch (error) {
                console.error('Error fetching custom system settings:', error);
            }
        };
        fetchCustomSystemSettings();
    }, []);

    const handlePublish = async (courseId) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/publish`, {
                method: 'PUT',
            });
            if (response.ok) {
                setCourses((prevCourses) =>
                    prevCourses.map((course) =>
                        course.id === courseId ? { ...course, published: true, available: true } : course
                    )
                );
                console.log('Course published successfully');
            } else {
                console.error('Failed to publish course');
            }
        } catch (error) {
            console.error('Error publishing course:', error);
        }
    };

    const handleAvailability = async (courseId, available) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/availability`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ available }),
            });
            if (response.ok) {
                setCourses((prevCourses) =>
                    prevCourses.map((course) =>
                        course.id === courseId ? { ...course, available } : course
                    )
                );
                console.log(`Course ${available ? 'set as available' : 'set as unavailable'} successfully`);
            } else {
                console.error(`Failed to ${available ? 'set course as available' : 'set course as unavailable'}`);
            }
        } catch (error) {
            console.error(`Error ${available ? 'setting course as available' : 'setting course as unavailable'}:`, error);
        }
    };

    if (!customSystemSettings) {
        return <div>Loading...</div>
    }

    return (
        <List>
            <Grid container spacing={4}>
                {courses.map((co, index) => {
                    const course = courses && courses.find(c => c.id === co.id);
                    return (course &&
                        <Grid item xs={12} sm={6} md={4} key={index}>
                            <Card>
                                <CardActionArea component={Link} to={`/coursebuilder/${course.id}/modulebuilder`} style={{ textDecoration: 'none' }}>
                                    <CardMedia
                                        component="img"
                                        height="140"
                                        image={course.icon}
                                        alt="Course image"
                                    />
                                    <CardContent>
                                        <Typography variant="h5" component="div">
                                            {course.title}
                                        </Typography>
                                        <div>
                                            <Chip
                                                label={course.published ? 'Published' : 'Unpublished'}
                                                color={course.published ? 'primary' : 'default'}
                                                size="small"
                                                style={{ marginRight: '8px' }}
                                            />
                                            <Chip
                                                label={course.available ? 'Available' : 'Unavailable'}
                                                color={course.available ? 'success' : 'default'}
                                                size="small"
                                            />
                                        </div>
                                    </CardContent>
                                </CardActionArea>
                                <CardActions>
                                    {!course.published && (
                                        <Button size="small" color="primary" onClick={() => handlePublish(course.id)}>
                                            Publish
                                        </Button>
                                    )}
                                    {course.published && (
                                        <Button
                                            size="small"
                                            color={course.available ? 'warning' : 'success'}
                                            onClick={() => handleAvailability(course.id, !course.available)}
                                        >
                                            {course.available ? 'Set Unavailable' : 'Set Available'}
                                        </Button>
                                    )}
                                    <Button size="small" color="secondary" component={Link} to={`/courses/${course.id}`}>
                                        Preview
                                    </Button>
                                </CardActions>
                            </Card>
                        </Grid>
                    );
                })}
            </Grid>
            {courses.length === 0 && (
                <Typography variant="body1">
                    You haven't created any courses yet.
                </Typography>
            )}
        </List>
    );
};

export default InstructorCourseList;