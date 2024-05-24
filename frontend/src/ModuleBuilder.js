import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, TextField, Button, IconButton } from '@mui/material';
import { Delete as DeleteIcon, DragIndicator as DragIndicatorIcon, Edit as EditIcon } from '@mui/icons-material';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import Layout from './Layout';
import { useParams, Link } from 'react-router-dom';

const ModuleBuilder = ({ user }) => {
    const [modules, setModules] = useState([]);
    const [moduleTitle, setModuleTitle] = useState('');
    const [moduleDescription, setModuleDescription] = useState('');
    const { courseId } = useParams();
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
        const fetchModules = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/modules`);
                if (response.ok) {
                    const data = await response.json();
                    setModules(data);
                } else {
                    console.error('Failed to fetch modules');
                }
            } catch (error) {
                console.error('Error fetching modules:', error);
            }
        };
        fetchModules();
    }, [courseId]);

    const handleModuleTitleChange = (event) => {
        setModuleTitle(event.target.value);
    };

    const handleModuleDescriptionChange = (event) => {
        setModuleDescription(event.target.value);
    };

    const handleAddModule = async () => {
        if (moduleTitle.trim() !== '') {
            try {
                const newModuleOrder = modules ? modules.length > 0 ? Math.max(...modules.map(module => module.order)) + 1 : 0 : 0;

                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/modules`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        title: moduleTitle,
                        description: moduleDescription,
                        order: newModuleOrder,
                    }),
                });

                if (response.ok) {
                    const newModule = await response.json();
                    if (modules) {
                        setModules([...modules, newModule]);
                    } else {
                        setModules([newModule]);
                    }
                    setModuleTitle('');
                } else {
                    console.error('Failed to add module');
                }
            } catch (error) {
                console.error('Error adding module:', error);
            }
        }
    };

    const handleDeleteModule = async (id) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/module/${id}`, {
                method: 'DELETE',
            });
            if (response.ok) {
                const updatedModules = modules.filter((module) => module.id !== id);
                setModules(updatedModules);
            } else {
                console.error('Failed to delete module');
            }
        } catch (error) {
            console.error('Error deleting module:', error);
        }
    };

    const onDragEnd = async (result) => {
        if (!result.destination) return;

        const newModules = Array.from(modules);
        const [reorderedModule] = newModules.splice(result.source.index, 1);
        newModules.splice(result.destination.index, 0, reorderedModule);

        try {
            const updatedOrder = newModules.reduce((acc, module, index) => {
                acc[module.id] = index;
                return acc;
            }, {});

            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/modules/order`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ order: updatedOrder }),
            });
            if (response.ok) {
                setModules(newModules);
            } else {
                console.error('Failed to update module order');
            }
        } catch (error) {
            console.error('Error updating module order:', error);
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
                        Module Builder
                    </Typography>
                    <Box mb={2}>
                        <TextField
                            label="Module Title"
                            value={moduleTitle}
                            onChange={handleModuleTitleChange}
                            fullWidth
                            margin="normal"
                        />

                        <TextField
                            label="Module Description"
                            value={moduleDescription}
                            onChange={handleModuleDescriptionChange}
                            fullWidth
                            margin="normal"
                            multiline
                            rows={4}
                        />

                        <Button variant="contained" color="primary" onClick={handleAddModule} style={{ marginLeft: '16px' }}>
                            Add
                        </Button>
                    </Box>
                </Box>
                <DragDropContext onDragEnd={onDragEnd}>
                    <Droppable droppableId="modules">
                        {(provided) => (
                            <Box {...provided.droppableProps} ref={provided.innerRef}>
                                {modules && modules.map((module, index) => (
                                    <Draggable key={module.id} draggableId={module.id.toString()} index={index}>
                                        {(provided) => (
                                            <Box
                                                ref={provided.innerRef}
                                                {...provided.draggableProps}
                                                display="flex"
                                                alignItems="center"
                                                justifyContent="space-between"
                                                mb={2}
                                            >
                                                <Box>
                                                    <Typography variant="h6">{module.title}</Typography>
                                                    <Typography variant="body2" color="textSecondary">{module.description}</Typography>
                                                </Box>
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
                                                        to={`/coursebuilder/${courseId}/modulebuilder/${module.id}/submodulebuilder`}
                                                        aria-label="edit"
                                                    >
                                                        <EditIcon />
                                                    </IconButton>
                                                    <IconButton
                                                        edge="end"
                                                        color="secondary"
                                                        onClick={() => handleDeleteModule(module.id)}
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

export default ModuleBuilder;
