import React, { useState, useEffect } from 'react';
import LoggedIn from './LoggedIn';
import NonLoggedIn from './NonLoggedIn';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';
import CoursePage from './CoursePage';
import { jwtDecode } from 'jwt-decode';
import ModulePage from './ModulePage';
import SubmodulePage from './SubmodulePage';
import CourseBuilder from './CourseBuilder';
import SubmoduleBuilder from './SubmoduleBuilder';
import ModuleBuilder from './ModuleBuilder';
import ElementBuilder from './ElementBuilder';
import { ThemeProvider, createTheme } from '@mui/material/styles';

function App() {
  const [token, setToken] = useState('');
  const [user, setUser] = useState(null);
  const [courses, setCourses] = useState(null);
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
    const storedToken = localStorage.getItem('token');
    setToken(storedToken);
  }, []);

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (token) {
      const decodedToken = jwtDecode(token);
      const { id, email, full_name } = decodedToken;

      setUser({
        id: id,
        name: full_name,
        email: email,
        avatar: 'https://via.placeholder.com/150',
        courses: [],
      });

      const fetchCourses = async () => {
        try {

          if (token) {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course`);
            if (response.ok) {
              const courses = await response.json();
              setCourses(courses);
            } else {
              console.error('Failed to fetch courses');
            }
          }
        } catch (error) {
          console.error('Error fetching courses:', error);
        }
      };

      const fetchUserData = async (userId) => {
        try {

          if (token) {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/user/` + userId, {
              headers: {
                'Authorization': `Bearer ${token}`,
              },
            });
            if (response.ok) {
              const userData = await response.json();

              // Update the user state with the fetched courses
              setUser((prevUser) => ({
                ...prevUser,
                ...userData,
              }));
            } else {
              console.error('Failed to fetch userData');
              localStorage.removeItem('token');
              window.location.href = '/';
            }
          }
        } catch (error) {
          console.error('Error fetching userData:', error);
          localStorage.removeItem('token');
          window.location.href = '/';
        }
      };

      fetchCourses();
      fetchUserData(id);
    }
  }, []);

  if (!customSystemSettings) {
    return <div>Loading...</div>;
  }

  document.title = customSystemSettings.title;

  return (
    <ThemeProvider theme={createTheme({
      palette: {
        primary: { main: customSystemSettings.primary_color },
        secondary: { main: customSystemSettings.secondary_color },
      }
    })}>
      <Router>
        <div className="App">
          <Routes>
            <Route path="/" element={
              <>
                {token ? (
                  <LoggedIn user={user} setUser={setUser} courses={courses} />
                ) : (
                  <NonLoggedIn />
                )}
              </>
            } />
            <Route path="/courses/:courseId" element={user && <CoursePage user={user} />} />
            <Route path="/courses/:courseId/modules/:moduleId" element={user && <ModulePage user={user} />} />
            <Route path="/courses/:courseId/modules/:moduleId/submodules/:submoduleId" element={user && <SubmodulePage user={user} />} />
            <Route path="/coursebuilder" element={user && <CourseBuilder user={user} />} />
            <Route path="/coursebuilder/:courseId/modulebuilder" element={user && <ModuleBuilder user={user} />} />
            <Route path="/coursebuilder/:courseId/modulebuilder/:moduleId/submodulebuilder" element={user && <SubmoduleBuilder user={user} />} />
            <Route path="/coursebuilder/:courseId/modulebuilder/:moduleId/submodulebuilder/:submoduleId/elementbuilder" element={user && <ElementBuilder user={user} />} />
          </Routes>
        </div>
      </Router>
    </ThemeProvider>
  );
}

export default App;