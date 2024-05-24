import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, Grid, Card, CardContent, Button, Breadcrumbs, LinearProgress } from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import Layout from './Layout';
import { useParams, Link } from 'react-router-dom';
import { Link as RouterLink } from 'react-router-dom';
import { Link as MuiLink } from '@mui/material';

const CoursePage = ({ user }) => {
    const { courseId } = useParams();
    const [course, setCourse] = useState(null);
    const [modules, setModules] = useState([]);
    const [completedModules, setCompletedModules] = useState([]);
    const [customSystemSettings, setCustomSystemSettings] = useState(null);

    const LinkRouter = (props) => <MuiLink {...props} underline="hover" sx={{ display: 'flex', alignItems: 'center' }} color="inherit" component={RouterLink} />;

    useEffect(() => {
        const fetchCourse = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}`);
                if (response.ok) {
                    const courseData = await response.json();
                    setCourse(courseData);
                } else {
                    console.error('Failed to fetch course');
                }
            } catch (error) {
                console.error('Error fetching course:', error);
            }
        };

        const fetchModules = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/modules`);
                if (response.ok) {
                    const modulesData = await response.json();
                    setModules(modulesData);
                } else {
                    console.error('Failed to fetch modules');
                }
            } catch (error) {
                console.error('Error fetching modules:', error);
            }
        };

        const fetchCompletedModules = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/courseProgress/${courseId}`);
                if (response.ok) {
                    const completedModulesData = await response.json();
                    setCompletedModules(completedModulesData);
                } else {
                    console.error('Failed to fetch completed modules');
                }
            } catch (error) {
                console.error('Error fetching completed modules:', error);
            }
        };

        fetchCourse();
        fetchModules();
        fetchCompletedModules();
    }, [courseId, user.id]);

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

    const isModuleCompleted = (moduleId) => {
        return completedModules && completedModules.find((m) => m.module_id === moduleId)?.progress === 1;
    };

    const isModuleAccessible = (index) => {
        return index === 0 || isModuleCompleted(modules[index - 1].id);
    };

    if (!customSystemSettings) {
        return <div>Loading...</div>
    }

    return (
        <Layout user={user}>
            {course && (
                <Container maxWidth="lg" className="my-8">
                    <Breadcrumbs aria-label="breadcrumb" style={{ marginBottom: '16px' }}>
                        <LinkRouter to={`/`}>
                            <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
                            Home
                        </LinkRouter>
                        <Typography color="textPrimary">{course.title}</Typography>
                    </Breadcrumbs>

                    <Box mb={4}>
                        <Typography variant="h4" component="h1" gutterBottom>
                            {course.title}
                        </Typography>
                        <Typography variant="subtitle1" color="textSecondary" gutterBottom>
                            Instructor: {course.instructor_full_name}
                        </Typography>
                        <Typography variant="subtitle1" color="textSecondary" gutterBottom>
                            Reward: {course.xp_reward} XP
                        </Typography>
                        <Typography variant="body1" gutterBottom>
                            {course.description}
                        </Typography>
                    </Box>

                    <Grid container spacing={4}>
                        <Grid item xs={12}>
                            <Box mb={4}>
                                <Typography variant="h5" gutterBottom>
                                    Course Modules
                                </Typography>
                                <Card>
                                    <CardContent>
                                        {modules && modules.map((m, index) => (
                                            <Box key={m.id} sx={{ mb: 2 }}>
                                                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                                                    <Box>
                                                        <Typography variant="h6">
                                                            {isModuleCompleted(m.id) ? '✓' : `${index + 1}.`} {m.title}
                                                        </Typography>
                                                        <Typography variant="body2" color="textSecondary">
                                                            {m.description}
                                                        </Typography>
                                                    </Box>
                                                    <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                                        <Typography variant="body2" color="textSecondary" sx={{ mr: 2 }}>
                                                            Reward: {m.xp_reward} XP
                                                        </Typography>
                                                        <Button
                                                            variant="contained"
                                                            color="primary"
                                                            component={Link}
                                                            to={`/courses/${courseId}/modules/${m.id}`}
                                                            disabled={!isModuleAccessible(index)}
                                                            size="small"
                                                        >
                                                            {isModuleCompleted(m.id) ? 'Revisit' : completedModules && completedModules.find((mp) => mp.module_id === m.id)?.progress * 100 > 0 ? 'Continue' : 'Start'}
                                                        </Button>
                                                    </Box>
                                                </Box>
                                                <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                                                    <Box sx={{ width: '100%', mr: 1 }}>
                                                        <LinearProgress
                                                            variant="determinate"
                                                            value={
                                                                (completedModules && completedModules.find((mp) => mp.module_id === m.id)?.progress * 100) || 0
                                                            }
                                                        />
                                                    </Box>
                                                    <Box sx={{ minWidth: 35 }}>
                                                        <Typography variant="body2" color="textSecondary">
                                                            {`${Math.round(
                                                                (completedModules && completedModules.find((mp) => mp.module_id === m.id)?.progress * 100) || 0
                                                            )}%`}
                                                        </Typography>
                                                    </Box>
                                                </Box>
                                            </Box>
                                        ))}
                                    </CardContent>
                                </Card>
                            </Box>
                        </Grid>
                    </Grid>
                </Container>
            )}
        </Layout>
    );
};

export default CoursePage;
