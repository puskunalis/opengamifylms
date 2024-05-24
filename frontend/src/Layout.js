import React from 'react';
import { AppBar, Toolbar, IconButton, Typography, Avatar, Popover, Button, Box } from '@mui/material';
import { Add as AddIcon } from '@mui/icons-material';
import AvatarMenu from './AvatarMenu';
import stringAvatar from './StringAvatar';
import { Link } from 'react-router-dom';

const Layout = ({ user, children }) => {
  const [anchorEl, setAnchorEl] = React.useState(null);
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

  const handleClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'simple-popover' : undefined;

  if (!customSystemSettings) {
      return <div>Loading...</div>
  }

  return (
    <Box>
      <AppBar position="static" sx={{ mb: 2 }}>
        <Toolbar>
          <IconButton edge="start" color="inherit" aria-label="menu">
            {/*<MenuIcon />*/}
          </IconButton>
          <Typography
            variant="h6"
            sx={{
              flexGrow: 1,
              textDecoration: 'none',
              color: 'white',
              fontWeight: 'bold',
              '&:hover': {
                color: 'rgba(255, 255, 255, 0.7)',
              },
            }}
            component={Link}
            to="/"
          >
            {customSystemSettings.title}
          </Typography>
          {user.is_instructor && <Button
            variant="contained"
            color="secondary"
            startIcon={<AddIcon />}
            component={Link}
            to="/coursebuilder"
            sx={{
              marginRight: 2,
            }}
          >
            Course Builder
          </Button>}
          <Avatar {...stringAvatar(user.name)} onClick={handleClick} />
        </Toolbar>
      </AppBar>
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
        slotProps={{
          paper: {
            style: {
              maxWidth: '600px',
            },
          },
        }}
      >
        <AvatarMenu />
      </Popover>
      {children}
    </Box>
  );
};

export default Layout;