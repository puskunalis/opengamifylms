import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, Grid, LinearProgress, Avatar, Card, CardContent, Tooltip } from '@mui/material';
import { styled } from '@mui/system';
import CourseDialog from './CourseDialog';
import getQuoteOfTheDay from './Quote';
import CourseList from './CourseList';
import UserCourseList from './UserCourseList';
import Layout from './Layout';

const LevelProgress = styled(LinearProgress)(({ theme }) => ({
    height: 20,
    borderRadius: 10,
}));

const ChallengeCard = styled(Card)(({ theme }) => ({
    transition: 'transform 0.3s',
    '&:hover': {
        transform: 'scale(1.05)',
    },
}));

const LoggedIn = ({ user, setUser, courses }) => {
    const [openDialog, setOpenDialog] = useState(false);
    const [selectedCourse, setSelectedCourse] = useState(null);
    const [userCourses, setUserCourses] = useState([]);
    const [leaderboardData, setLeaderboardData] = useState([]);
    const [userBadges, setUserBadges] = useState([]);

    useEffect(() => {
        const fetchUserData = async () => {
            try {
                const token = localStorage.getItem('token');

                // Fetch user courses
                const userCoursesResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/enrollments`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                if (userCoursesResponse.ok) {
                    const userCoursesData = await userCoursesResponse.json();
                    setUserCourses(userCoursesData ? userCoursesData : []);
                } else {
                    console.error('Failed to fetch user courses');
                }

                // Fetch leaderboard users
                const leaderboardResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/topUsersByXp`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                if (leaderboardResponse.ok) {
                    const leaderboardData = await leaderboardResponse.json();
                    setLeaderboardData(leaderboardData ? leaderboardData : []);
                } else {
                    console.error('Failed to fetch leaderboard users');
                }

                // Fetch user badges
                const userBadgesResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/badges`, {
                    headers: {
                        'Authorization': `Bearer ${token}`,
                    },
                });
                if (userBadgesResponse.ok) {
                    const userBadgesData = await userBadgesResponse.json();
                    setUserBadges(userBadgesData ? userBadgesData : []);
                } else {
                    console.error('Failed to fetch user badges');
                }
            } catch (error) {
                console.error('Error fetching user courses:', error);
            }
        };

        if (user) {
            fetchUserData();
        }
    }, [user]);

    const handleEnrollment = async (user, course) => {
        try {
            const token = localStorage.getItem('token');
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/enrollment`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`,
                },
                body: JSON.stringify({ userId: user.id, courseId: course.id }),
            });

            if (response.ok) {
                // Enrollment successful, update the user's courses
                setOpenDialog(false);
                setUser((prevUser) => {
                    const existingCourse = prevUser.courses?.find(c => c.id === course.id);
                    if (!existingCourse) {
                        return {
                            ...prevUser,
                            courses: prevUser.courses ? [...prevUser.courses, { id: course.id, progress: 0 }] : [{ id: course.id, progress: 0 }],
                        };
                    }
                    return prevUser;
                });
            } else {
                console.error('Failed to enroll in course');
            }
        } catch (error) {
            console.error('Error enrolling in course:', error);
        }
    };

    if (!user) {
        return <div>Loading...</div>;
    }

    return (
        <Layout user={user}>
            <Box mt={4}>
                <Container maxWidth="lg">
                    <Grid container spacing={4}>
                        <Grid item xs={12}>
                            <Box display="flex" alignItems="center" mb={2}>
                                <Box>
                                    <Typography variant="h4" component="h1">Welcome, {user.name}!</Typography>
                                    <Typography variant="subtitle1" sx={{ fontStyle: 'italic' }}>
                                        {getQuoteOfTheDay()}
                                    </Typography>
                                </Box>
                            </Box>
                            <Box mb={4}>
                                <Typography variant="h5" gutterBottom>Your Level: {1 + Math.floor(user.xp / 100)}</Typography>
                                <LevelProgress variant="determinate" value={(user.xp % 100) / 100 * 100} />
                                <Typography variant="body1">{user.xp % 100} / 100 XP</Typography>
                            </Box>
                            <Box mb={4}>
                                <Typography variant="h5" gutterBottom>Challenges</Typography>
                                <Grid container spacing={2}>
                                    <Grid item xs={12} sm={6} md={4}>
                                        <ChallengeCard>
                                            <CardContent>
                                                <Typography variant="h6">Complete 5 Courses</Typography>
                                                <Typography variant="body2" color="primary">Reward: 100 XP</Typography>
                                                <LinearProgress variant="determinate" value={(2 /* TODO */ / 5) * 100} title="2/5 Courses Completed" />
                                            </CardContent>
                                        </ChallengeCard>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={4}>
                                        <ChallengeCard>
                                            <CardContent>
                                                <Typography variant="h6">Earn 500 XP</Typography>
                                                <Typography variant="body2" color="primary">Reward: 50 XP</Typography>
                                                <LinearProgress variant="determinate" value={(user.xp / 500) * 100} title={`${user.xp}/500 XP Earned`} />
                                            </CardContent>
                                        </ChallengeCard>
                                    </Grid>
                                    <Grid item xs={12} sm={6} md={4}>
                                        <ChallengeCard>
                                            <CardContent>
                                                <Typography variant="h6">Enroll in 3 Courses</Typography>
                                                <Typography variant="body2" color="primary">Reward: 75 XP</Typography>
                                                <LinearProgress variant="determinate" value={(userCourses.length / 3) * 100} title={`${userCourses.length}/3 Courses Enrolled`} />
                                            </CardContent>
                                        </ChallengeCard>
                                    </Grid>
                                </Grid>
                            </Box>
                        </Grid>
                        <Grid item xs={12} md={8}>
                            <UserCourseList courses={userCourses} user={user} />
                            <Box mt={4}>
                                {courses && userCourses && (
                                    <CourseList
                                        courses={courses.filter(
                                            (course) => !userCourses.some((userCourse) => userCourse.id === course.id)
                                        )}
                                        setSelectedCourse={setSelectedCourse}
                                        setOpenDialog={setOpenDialog}
                                        user={user}
                                    />
                                )}
                            </Box>
                        </Grid>
                        <Grid item xs={12} md={4}>
                            <Typography variant="h4" gutterBottom>Leaderboard</Typography>
                            <Card>
                                <CardContent>
                                    {leaderboardData.map((entry, index) => (
                                        <Box key={entry.id} display="flex" alignItems="center" justifyContent="space-between" mb={1}>
                                            <Typography variant="h6" mr={2}>{index + 1}. {entry.full_name}</Typography>
                                            <Typography variant="subtitle1">{entry.xp} XP</Typography>
                                        </Box>
                                    ))}
                                </CardContent>
                            </Card>
                            <Box mt={4}>
                                <Typography variant="h4" gutterBottom>Badges</Typography>
                                <Box display="flex" flexWrap="wrap">
                                    {userBadges.map((badge) => (
                                        <Tooltip key={badge.id} title={badge.description} placement="bottom">
                                            <Box display="flex" flexDirection="column" alignItems="center" mr={2}>
                                                <Avatar variant="rounded" src={badge.icon} sx={{ width: 100, height: 100, borderRadius: '10%' }} />
                                                <Typography variant="subtitle1">{badge.title}</Typography>
                                            </Box>
                                        </Tooltip>
                                    ))}
                                </Box>
                            </Box>
                        </Grid>
                    </Grid>
                </Container>
            </Box>
            {selectedCourse && (
                <CourseDialog open={openDialog} handleClose={() => setOpenDialog(false)} user={user} course={selectedCourse} onEnroll={handleEnrollment} />
            )}
        </Layout>
    );
};

export default LoggedIn;
