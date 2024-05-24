import React from 'react';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogActions from '@mui/material/DialogActions';
import Button from '@mui/material/Button';

function CourseDialog({ open, handleClose, user, course, onEnroll }) {
 return (
    <Dialog open={open} onClose={handleClose}>
      <DialogTitle>{course?.title}</DialogTitle>
      <DialogContent>
        <DialogContentText>
          Description: {course.description}
          <br />
          Reward: {course.xp_reward} XP
          <br />
          Instructor: {course.instructor_full_name}
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose}>Close</Button>
        <Button onClick={() => onEnroll(user, course)}>Enroll</Button>
      </DialogActions>
    </Dialog>
 );
}

export default CourseDialog;