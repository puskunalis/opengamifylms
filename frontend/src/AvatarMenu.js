import * as React from 'react';
import List from '@mui/material/List';
import ListItemText from '@mui/material/ListItemText';
import ListItemButton from '@mui/material/ListItemButton';

export default function PinnedSubheaderList() {
    return (
        <List
            sx={{
                width: '100%',
                maxWidth: 360,
                bgcolor: 'background.paper',
                position: 'relative',
                overflow: 'auto',
                maxHeight: 300,
                '& ul': { padding: 0 },
            }}
            subheader={<li />}
        >
            <ListItemButton key={`item-3`} onClick={() => {
                localStorage.removeItem('token');
                window.location.href = '/';
            }}>
                <ListItemText primary={`Log out`} />
            </ListItemButton>
        </List>
    );
}