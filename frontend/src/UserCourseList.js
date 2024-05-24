import { Typography, Grid, Card, CardContent, CardMedia, CardActionArea, Box, LinearProgress } from '@mui/material';
import { Link } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { styled } from '@mui/system';

const CourseCard = styled(Card)(({ theme }) => ({
    transition: 'transform 0.3s',
    '&:hover': {
        transform: 'scale(1.05)',
    },
}));

const XpProgress = styled(LinearProgress)(({ theme }) => ({
    height: 20,
    borderRadius: 10,
}));

const UserCourseList = ({ courses, user }) => {
    const [courseProgress, setCourseProgress] = useState({});
    const [customSystemSettings, setCustomSystemSettings] = useState(null);

    useEffect(() => {
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

    useEffect(() => {
        const fetchCourseProgress = async () => {
            try {
                const promises = courses.map(async (course) => {
                    const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/courseProgress/${course.id}`);
                    if (response.ok) {
                        const completedModulesData = await response.json();
                        return { courseId: course.id, completedModules: completedModulesData };
                    } else {
                        console.error('Failed to fetch completed modules');
                    }
                });

                const results = await Promise.all(promises);
                const progress = {};
                results.forEach((result) => {
                    if (result) {
                        progress[result.courseId] = result.completedModules;
                    }
                });
                setCourseProgress(progress);
            } catch (error) {
                console.error('Error fetching course progress:', error);
            }
        };

        fetchCourseProgress();
    }, [courses, user.id]);

    const calculateCourseProgress = (courseId) => {
        const completedModules = courseProgress[courseId] || [];
        const totalModules = completedModules.length;
        const totalProgress = completedModules.reduce((sum, module) => sum + module.progress, 0);
        return totalModules > 0 ? Math.round((totalProgress / totalModules) * 100) : 100;
    };

    if (!customSystemSettings) {
        return <div>Loading...</div>
    }

    return (
        <Box>
            <Typography variant="h4" component="h1" gutterBottom>
                Your Adventure
            </Typography>
            <Grid container spacing={4}>
                {courses.map((course, index) => (
                    <Grid item xs={12} sm={6} md={4} key={index} component={Link} to={`/courses/${course.id}`} style={{ textDecoration: 'none' }}>
                        <CourseCard>
                            <CardActionArea>
                                <CardMedia component="img" height="200" image={course.icon} alt="Course image" />
                                <CardContent>
                                    <Typography variant="h5" component="div" gutterBottom>
                                        {course.title}
                                    </Typography>
                                    <Box display="flex" alignItems="center" mb={1}>
                                        <Typography variant="body1" color="text.secondary" mr={1}>
                                            Progress:
                                        </Typography>
                        <Typography variant="body1" color="primary">
                                            {calculateCourseProgress(course.id)}%
                                        </Typography>
                                    </Box>
                                    <XpProgress variant="determinate" value={calculateCourseProgress(course.id)} />
                                </CardContent>
                            </CardActionArea>
                        </CourseCard>
                    </Grid>
                ))}
            </Grid>
        </Box>
    );
};


export default UserCourseList;
