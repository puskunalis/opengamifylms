import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, TextField, Button, IconButton } from '@mui/material';
import { Delete as DeleteIcon, DragIndicator as DragIndicatorIcon, Edit as EditIcon } from '@mui/icons-material';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import Layout from './Layout';
import { useParams, Link } from 'react-router-dom';

const SubmoduleBuilder = ({ user }) => {
    const [submodules, setSubmodules] = useState([]);
    const [submoduleTitle, setSubmoduleTitle] = useState('');
    const { courseId, moduleId } = useParams();
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
        const fetchSubmodules = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/module/${moduleId}/submodules`);
                if (response.ok) {
                    const data = await response.json();
                    setSubmodules(data);
                } else {
                    console.error('Failed to fetch submodules');
                }
            } catch (error) {
                console.error('Error fetching submodules:', error);
            }
        };
        fetchSubmodules();
    }, [courseId, moduleId]);

    const handleSubmoduleTitleChange = (event) => {
        setSubmoduleTitle(event.target.value);
    };

    const handleAddSubmodule = async () => {
        if (submoduleTitle.trim() !== '') {
            try {
                const newSubmoduleOrder = submodules ? submodules.length > 0 ? Math.max(...submodules.map(submodule => submodule.order)) + 1 : 0 : 0;

                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/module/${moduleId}/submodules`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        title: submoduleTitle,
                        order: newSubmoduleOrder,
                    }),
                });

                if (response.ok) {
                    const newSubmodule = await response.json();
                    if (submodules) {
                        setSubmodules([...submodules, newSubmodule]);
                    } else {
                        setSubmodules([newSubmodule]);
                    }
                    setSubmoduleTitle('');
                } else {
                    console.error('Failed to add submodule');
                }
            } catch (error) {
                console.error('Error adding submodule:', error);
            }
        }
    };

    const handleDeleteSubmodule = async (id) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/submodule/${id}`, {
                method: 'DELETE',
            });
            if (response.ok) {
                const updatedSubmodules = submodules.filter((submodule) => submodule.id !== id);
                setSubmodules(updatedSubmodules);
            } else {
                console.error('Failed to delete submodule');
            }
        } catch (error) {
            console.error('Error deleting submodule:', error);
        }
    };

    const onDragEnd = async (result) => {
        if (!result.destination) return;

        const newSubmodules = Array.from(submodules);
        const [reorderedSubmodule] = newSubmodules.splice(result.source.index, 1);
        newSubmodules.splice(result.destination.index, 0, reorderedSubmodule);

        try {
            const updatedOrder = newSubmodules.reduce((acc, submodule, index) => {
                acc[submodule.id] = index;
                return acc;
            }, {});

            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/module/${moduleId}/submodules/order`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ order: updatedOrder }),
            });
            if (response.ok) {
                setSubmodules(newSubmodules);
            } else {
                console.error('Failed to update submodule order');
            }
        } catch (error) {
            console.error('Error updating submodule order:', error);
        }
    };

    if (!customSystemSettings) {
        return <div>Loading...</div>
    }

    return (
        <Layout user={user}>
            <Container maxWidth="lg" className="my-8">
                <Box mb={4}>
                    <Typography variant="h4" component="h1" gutterBottom>
                        Submodule Builder
                    </Typography>
                    <Box display="flex" alignItems="center" mb={2}>
                        <TextField
                            label="Submodule Title"
                            value={submoduleTitle}
                            onChange={handleSubmoduleTitleChange}
                            fullWidth
                            margin="normal"
                        />

                        <Button variant="contained" color="primary" onClick={handleAddSubmodule} style={{ marginLeft: '16px' }}>
                            Add
                        </Button>
                    </Box>
                </Box>
                <DragDropContext onDragEnd={onDragEnd}>
                    <Droppable droppableId="submodules">
                        {(provided) => (
                            <Box {...provided.droppableProps} ref={provided.innerRef}>
                                {submodules && submodules.map((submodule, index) => (
                                    <Draggable key={submodule.id} draggableId={submodule.id.toString()} index={index}>
                                        {(provided) => (
                                            <Box
                                                ref={provided.innerRef}
                                                {...provided.draggableProps}
                                                display="flex"
                                                alignItems="center"
                                                justifyContent="space-between"
                                                mb={2}
                                            >
                                                <Typography variant="h6">{submodule.title}</Typography>
                                                <Box display="flex">
                                                    <IconButton
                                                        {...provided.dragHandleProps}
                                                        color="primary"
                                                        aria-label="drag"
                                                    >
                                                        <DragIndicatorIcon />
                                                    </IconButton>
                                                    <IconButton
                                                        color="primary"
                                                        component={Link}
                                                        to={`/coursebuilder/${courseId}/modulebuilder/${moduleId}/submodulebuilder/${submodule.id}/elementbuilder`}
                                                        aria-label="edit"
                                                    >
                                                        <EditIcon />
                                                    </IconButton>
                                                    <IconButton
                                                        edge="end"
                                                        color="secondary"
                                                        onClick={() => handleDeleteSubmodule(submodule.id)}
                                                        aria-label="delete"
                                                    >
                                                        <DeleteIcon />
                                                    </IconButton>
                                                </Box>
                                            </Box>
                                        )}
                                    </Draggable>
                                ))}
                                {provided.placeholder}
                            </Box>
                        )}
                    </Droppable>
                </DragDropContext>
            </Container>
        </Layout>
    );
};

export default SubmoduleBuilder;