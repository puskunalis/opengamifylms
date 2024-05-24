import { Typography, Grid, Card, CardContent, CardMedia, CardActionArea } from '@mui/material';

const CourseList = ({ courses, setSelectedCourse, setOpenDialog, user }) => (
  <>
    {user && courses && (
      <>
        {(user.courses === null || courses.filter((course) => !(user.courses || []).some(userCourse => userCourse.id === course.id)).length > 0) && (
          <Typography variant="h4" component="h1" gutterBottom>
            New Quests
          </Typography>
        )}
        <Grid container spacing={4}>
          {courses.filter((course) => !(user.courses || []).some(userCourse => userCourse.id === course.id)).map((course, index) => (
            <Grid item xs={12} sm={6} md={4} key={index} onClick={() => {
              setSelectedCourse(course);
              setOpenDialog(true);
            }}>
              <Card>
                <CardActionArea>
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
                  </CardContent>
                </CardActionArea>
              </Card>
            </Grid>
          ))}
        </Grid>
      </>
    )}
  </>
);

export default CourseList;