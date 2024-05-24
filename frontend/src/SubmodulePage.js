import React, { useState, useEffect } from 'react';
import { Container, Typography, Grid, Button, Box, Tooltip, IconButton, Breadcrumbs } from '@mui/material';
import { ArrowBack, ArrowForward, FiberManualRecord } from '@mui/icons-material';
import HomeIcon from '@mui/icons-material/Home';
import { useNavigate, useParams } from 'react-router-dom';
import { CSSTransition, TransitionGroup } from 'react-transition-group';
import Layout from './Layout';
import Quiz from './Quiz';
import './SubmodulePage.css';
import { Link as RouterLink } from 'react-router-dom';
import { Link as MuiLink } from '@mui/material';

const SubmodulePage = ({ user }) => {
    const { courseId, moduleId, submoduleId } = useParams();
    const [submodule, setSubmodule] = useState(null);
    const [submodules, setSubmodules] = useState(null);
    const [course, setCourse] = useState(null);
    const [module, setModule] = useState(null);
    const navigate = useNavigate();
    const [quizzesCompleted, setQuizzesCompleted] = useState(false);
    const [quizStates, setQuizStates] = useState({});
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

    const handleQuizComplete = (quizId, completed) => {
        setQuizStates((prevState) => ({
            ...prevState,
            [quizId]: completed,
        }));
    };

    useEffect(() => {
        const checkQuizzesCompleted = () => {
            if (submodule && submodule.elements) {
                const quizzes = submodule.elements.filter(
                    (el) => el.type === "quiz_single_choice" || el.type === "quiz_multiple_choice"
                );
                const allQuizzesCompleted = quizzes.every((quiz) => quizStates[quiz.quiz_id]);
                setQuizzesCompleted(allQuizzesCompleted);
            }
        };

        checkQuizzesCompleted();
    }, [submodule, quizStates, submoduleId]);

    const LinkRouter = (props) => <MuiLink {...props} underline="hover" sx={{ display: 'flex', alignItems: 'center' }} color="inherit" component={RouterLink} />;

    useEffect(() => {
        const fetchSubmodule = async () => {
            try {
                // Fetch the submodule data
                const submoduleResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/submodule/${submoduleId}`);
                if (submoduleResponse.ok) {
                    const fetchedSubmodule = await submoduleResponse.json();

                    // Fetch the elements for the submodule
                    const elementsResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/submodule/${submoduleId}/elements`);
                    if (elementsResponse.ok) {
                        const elements = await elementsResponse.json();
                        setSubmodule({ ...fetchedSubmodule, elements });
                    } else {
                        console.error('Failed to fetch elements');
                    }
                } else {
                    console.error('Failed to fetch submodule');
                }

                const submodulesResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/module/${moduleId}/submodules`);
                if (submodulesResponse.ok) {
                    const fetchedSubmodules = await submodulesResponse.json();
                    setSubmodules(fetchedSubmodules);
                } else {
                    console.error('Failed to fetch submodules');
                }
            } catch (error) {
                console.error('Error fetching submodule:', error);
            }
        };

        fetchSubmodule();

        const fetchCourse = async () => {
            try {
                const courseResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}`);
                if (courseResponse.ok) {
                    const fetchedCourse = await courseResponse.json();
                    setCourse(fetchedCourse);
                } else {
                    console.error('Failed to fetch course');
                }
            } catch (error) {
                console.error('Error fetching course:', error);
            }
        };

        const fetchModule = async () => {
            try {
                const moduleResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/module/${moduleId}`);
                if (moduleResponse.ok) {
                    const fetchedModule = await moduleResponse.json();
                    setModule(fetchedModule);
                } else {
                    console.error('Failed to fetch module');
                }
            } catch (error) {
                console.error('Error fetching module:', error);
            }
        };

        fetchCourse();
        fetchModule();
    }, [courseId, moduleId, submoduleId]);

    const handlePrevSubmodule = () => {
        const currentIndex = submodules.findIndex(s => s.id === parseInt(submoduleId));
        if (currentIndex > 0) {
            const prevSubmoduleId = submodules[currentIndex - 1].id;
            navigate(`/courses/${courseId}/modules/${moduleId}/submodules/${prevSubmoduleId}`);
        }
    };

    const handleNextSubmodule = async () => {
        const currentIndex = submodules.findIndex(s => s.id === parseInt(submoduleId));
        if (currentIndex < submodules.length - 1) {
            // Mark the current submodule as completed
            await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/completedSubmodule/${submoduleId}`, {
                method: 'PUT',
            });

            const nextSubmoduleId = submodules[currentIndex + 1].id;
            navigate(`/courses/${courseId}/modules/${moduleId}/submodules/${nextSubmoduleId}`);
        }
    };

    const handleFinishModule = async () => {
        // Mark the current submodule as completed
        await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/completedSubmodule/${submoduleId}`, {
            method: 'PUT',
        });

        // Redirect back to the course page
        navigate(`/courses/${courseId}`);
    };

    const getPrevSubmoduleTitle = () => {
        const currentIndex = submodules.findIndex(s => s.id === parseInt(submoduleId));
        if (currentIndex > 0) {
            return submodules[currentIndex - 1].title;
        }
        return '';
    };

    const getNextSubmoduleTitle = () => {
        const currentIndex = submodules.findIndex(s => s.id === parseInt(submoduleId));
        if (currentIndex < submodules.length - 1) {
            return submodules[currentIndex + 1].title;
        }
        return '';
    };

    if (!customSystemSettings) {
        return <div>Loading...</div>
    }

    return (
        <Layout user={user}>
            <TransitionGroup>
                {course && module && submodule && submodules && (<CSSTransition
                    key={submoduleId}
                    timeout={300}
                    classNames="fade"
                >
                    <Container maxWidth="lg" className="my-8">
                        <Breadcrumbs aria-label="breadcrumb" style={{ marginBottom: '16px' }}>
                            <LinkRouter to={`/`}>
                                <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
                                Home
                            </LinkRouter>
                            <LinkRouter to={`/courses/${courseId}`}>
                                {course.title}
                            </LinkRouter>
                            <LinkRouter to={`/courses/${courseId}/modules/${moduleId}`}>
                                {module.title}
                            </LinkRouter>
                            <Typography color="textPrimary">{submodule.title}</Typography>
                        </Breadcrumbs>
                        <Box mb={4}>
                            <Typography variant="h4" component="h1" gutterBottom>
                                {submodule.title}
                            </Typography>
                            {submodule.elements && submodule.elements.map((el) => {
                                if (el.type === "html") {
                                    return <div key={el.id} dangerouslySetInnerHTML={{ __html: el.content }} />;
                                } else if (el.type === "quiz_single_choice" || el.type === "quiz_multiple_choice") {
                                    return (
                                        <Quiz key={el.id} quizId={el.quiz_id} onQuizComplete={handleQuizComplete} />
                                    );
                                } else if (el.type === "video") {
                                    return (
                                        <video
                                            src={`/api/v1/element/video/${el.id}`}
                                            controls
                                            style={{
                                                width: '100%',
                                                height: 'auto',
                                                maxWidth: '100%',
                                            }}
                                            preload="auto"
                                            playsInline
                                        />
                                    );
                                }
                                return null;
                            })}
                        </Box>
                        <Grid container spacing={2} alignItems="center">
                            <Grid item xs={5} container justifyContent="left">
                                <Button
                                    variant="contained"
                                    startIcon={<ArrowBack />}
                                    onClick={handlePrevSubmodule}
                                    disabled={submodules.findIndex(s => s.id === parseInt(submoduleId)) === 0}
                                >
                                    Previous
                                </Button>
                                <Typography variant="subtitle1">{getPrevSubmoduleTitle()}</Typography>
                            </Grid>
                            <Grid item xs={2} container justifyContent="center">
                                {submodules.map((submodule, index) => (
                                    <Grid item key={submodule.id}>
                                        <Tooltip title={submodule.title} arrow>
                                            <IconButton
                                                size="small"
                                                onClick={() => navigate(`/courses/${courseId}/modules/${moduleId}/submodules/${submodule.id}`)}
                                                style={{
                                                    backgroundColor: submodule.id === parseInt(submoduleId) ? 'primary.main' : 'inherit',
                                                    color: submodule.id === parseInt(submoduleId) ? 'white' : 'black',
                                                    borderRadius: '50%',
                                                    border: submodule.id === parseInt(submoduleId) ? '1px solid black' : 'none',
                                                }}
                                            >
                                                <FiberManualRecord />
                                            </IconButton>
                                        </Tooltip>
                                    </Grid>
                                ))}
                            </Grid>
                            <Grid item xs={5} container justifyContent="right">
                                {submodules.findIndex((s) => s.id === parseInt(submoduleId)) === submodules.length - 1 ? (
                                    <Button
                                        variant="contained"
                                        onClick={handleFinishModule}
                                        disabled={!quizzesCompleted}
                                    >
                                        Finish Module
                                    </Button>
                                ) : (
                                    <>
                                        <Typography variant="subtitle1">{getNextSubmoduleTitle()}</Typography>
                                        <Button
                                            variant="contained"
                                            endIcon={<ArrowForward />}
                                            onClick={handleNextSubmodule}
                                            disabled={!quizzesCompleted || submodules.findIndex((s) => s.id === parseInt(submoduleId)) === submodules.length - 1}
                                        >
                                            Next
                                        </Button>
                                    </>
                                )}
                            </Grid>
                        </Grid>
                    </Container>
                </CSSTransition>)}
            </TransitionGroup>
        </Layout>
    );
}

export default SubmodulePage;
