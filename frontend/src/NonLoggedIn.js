import React, { useState } from 'react';
import { Button, Dialog, DialogTitle, DialogContent, TextField, Typography, Container, Box, Grid, Card } from '@mui/material';
import { styled } from '@mui/system';

const StyledContainer = styled(Container)(({ theme }) => ({
  backgroundImage: `url(https://picsum.photos/seed/nonloggedin14/3840/2160)`,
  backgroundSize: 'cover',
  backgroundPosition: 'center',
  minHeight: '100vh',
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'center',
  alignItems: 'center',
  padding: theme.spacing(4),
}));

const StyledCard = styled(Card)(({ theme }) => ({
  maxWidth: 400,
  margin: 'auto',
  padding: theme.spacing(4),
  boxShadow: '0px 4px 20px rgba(0, 0, 0, 0.1)',
  backgroundColor: 'rgba(255, 255, 255, 0.9)',
}));

function NonLoggedIn() {
  const [showLogin, setShowLogin] = useState(false);
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleLogin = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('token', data.token);
        window.location.reload();
      } else {
        setErrorMessage('Invalid email or password');
      }
    } catch (error) {
      console.error('Error during login:', error);
      setErrorMessage('An error occurred. Please try again.');
    }
  };

  return (
    <>
      <StyledContainer maxWidth={false}>
        <Grid container spacing={4} justifyContent="center" alignItems="center">
          <Grid item xs={12} md={6}>
            <Typography variant="h2" component="h1" gutterBottom color="textPrimary" align="center">
              Welcome to {document.title}
            </Typography>
            <Typography variant="h5" component="p" color="textSecondary" align="center" paragraph>
              Unlock the power of gamified learning and take your skills to the next level!
            </Typography>
            <Box display="flex" justifyContent="center" mt={4}>
              <Button variant="contained" color="primary" size="large" onClick={() => setShowLogin(true)}>
                Get Started
              </Button>
            </Box>
          </Grid>
        </Grid>
      </StyledContainer>

      <Dialog open={showLogin} onClose={() => setShowLogin(false)}>
        <StyledCard>
          <DialogTitle>Login</DialogTitle>
          <DialogContent>
            {errorMessage && (
              <Typography color="error" align="center" gutterBottom>
                {errorMessage}
              </Typography>
            )}
            <TextField
              autoFocus
              margin="dense"
              label="Email"
              type="email"
              fullWidth
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <TextField
              margin="dense"
              label="Password"
              type="password"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <Box mt={4} display="flex" justifyContent="center">
              <Button variant="contained" color="primary" onClick={handleLogin}>
                Login
              </Button>
            </Box>
          </DialogContent>
        </StyledCard>
      </Dialog>
    </>
  );
}

export default NonLoggedIn;
