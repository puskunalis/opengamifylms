import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, Breadcrumbs, List, ListItem, ListItemIcon, ListItemText, LinearProgress } from '@mui/material';
import HomeIcon from '@mui/icons-material/Home';
import LockIcon from '@mui/icons-material/Lock';
import LockOpenIcon from '@mui/icons-material/LockOpen';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import Layout from './Layout';
import { useParams, Link as RouterLink } from 'react-router-dom';
import { Link as MuiLink } from '@mui/material';

const ModulePage = ({ user }) => {
  const { courseId, moduleId } = useParams();
  const [module, setModule] = useState(null);
  const [submodules, setSubmodules] = useState([]);
  const [course, setCourse] = useState(null);
  const [completedSubmodules, setCompletedSubmodules] = useState([]);

  const LinkRouter = (props) => <MuiLink {...props} underline="hover" sx={{ display: 'flex', alignItems: 'center' }} color="inherit" component={RouterLink} />;

  useEffect(() => {
    const fetchModule = async () => {
      try {
        const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/module/${moduleId}`);
        if (response.ok) {
          const moduleData = await response.json();
          setModule(moduleData);
        } else {
          console.error('Failed to fetch module');
        }
      } catch (error) {
        console.error('Error fetching module:', error);
      }
    };

    const fetchSubmodules = async () => {
      try {
        const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/module/${moduleId}/submodules`);
        if (response.ok) {
          const submodulesData = await response.json();
          setSubmodules(submodulesData);
        } else {
          console.error('Failed to fetch submodules');
        }
      } catch (error) {
        console.error('Error fetching submodules:', error);
      }
    };

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

    const fetchCompletedSubmodules = async () => {
      try {
        const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/${user.id}/completedSubmodules`);
        if (response.ok) {
          const completedSubmodulesData = await response.json();
          setCompletedSubmodules(completedSubmodulesData);
        } else {
          console.error('Failed to fetch completed submodules');
        }
      } catch (error) {
        console.error('Error fetching completed submodules:', error);
      }
    };

    fetchModule();
    fetchSubmodules();
    fetchCourse();
    fetchCompletedSubmodules();
  }, [courseId, moduleId, user.id]);

  const isSubmoduleCompleted = (submoduleId) => {
    return completedSubmodules && completedSubmodules.some(completedSubmodule => completedSubmodule.submodule_id === submoduleId);
  };

  const isSubmoduleAccessible = (submoduleIndex) => {
    if (submoduleIndex === 0) {
      return true;
    }
    const previousSubmoduleId = submodules[submoduleIndex - 1].id;
    return isSubmoduleCompleted(previousSubmoduleId);
  };

  const getProgressPercentage = () => {
    if (!submodules) {
      return 100;
    }
    if (submodules.length === 0) {
      return 0;
    }
    const completedCount = submodules.filter(submodule => isSubmoduleCompleted(submodule.id)).length;
    return Math.round((completedCount / submodules.length) * 100);
  };

  return (
    <Layout user={user}>
      <Container maxWidth="lg" className="my-8">
        {course && module && (
          <Breadcrumbs aria-label="breadcrumb" style={{ marginBottom: '16px' }}>
            <LinkRouter to={`/`}>
              <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
              Home
            </LinkRouter>
            <LinkRouter to={`/courses/${courseId}`}>
              {course.title}
            </LinkRouter>
            <Typography color="textPrimary">{module.title}</Typography>
          </Breadcrumbs>
        )}
        <Box mb={4}>
          {module && (
            <Typography variant="h4" component="h1" gutterBottom>
              {module.title}
            </Typography>
          )}
          <Typography variant="body1" gutterBottom>
            Progress: {getProgressPercentage()}%
          </Typography>
          <LinearProgress variant="determinate" value={getProgressPercentage()} />
        </Box>
        <List>
          {submodules && submodules.map((submodule, index) => (
            <ListItem
              key={submodule.id}
              button
              component={RouterLink}
              to={isSubmoduleAccessible(index) ? `/courses/${courseId}/modules/${moduleId}/submodules/${submodule.id}` : '#'}
              disabled={!isSubmoduleAccessible(index)}
            >
              <ListItemIcon>
                {isSubmoduleCompleted(submodule.id) ? (
                  <CheckCircleIcon style={{ color: "success" }} />
                ) : isSubmoduleAccessible(index) ? (
                  <LockOpenIcon color="action" />
                ) : (
                  <LockIcon color="disabled" />
                )}
              </ListItemIcon>
              <ListItemText
                primary={submodule.title}
                secondary={`Reward: ${submodule.xp_reward} XP`}
              />
            </ListItem>
          ))}
        </List>
      </Container>
    </Layout>
  );
};

export default ModulePage;