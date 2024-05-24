import React, { useState, useEffect } from 'react';
import { Container, Typography, Box, TextField, Button, IconButton, FormLabel, FormControl, FormControlLabel, Radio, RadioGroup } from '@mui/material';
import { Delete as DeleteIcon, DragIndicator as DragIndicatorIcon } from '@mui/icons-material';
import { DragDropContext, Droppable, Draggable } from 'react-beautiful-dnd';
import Layout from './Layout';
import Quiz from './Quiz';
import { useParams } from 'react-router-dom';
import RichTextEditor from './RichTextEditor';
import { stateToHTML } from 'draft-js-export-html';

const ElementBuilder = ({ user }) => {
    const [htmlContent, setHtmlContent] = useState('');
    const [videoFile, setVideoFile] = useState(null);
    const [elements, setElements] = useState([]);
    const [quizQuestion, setQuizQuestion] = useState('');
    const [quizAnswers, setQuizAnswers] = useState(['']);
    const [quizCorrectAnswer, setQuizCorrectAnswer] = useState('');
    const [quizCorrectAnswers, setQuizCorrectAnswers] = useState([]);
    const [quizType, setQuizType] = useState('single');
    const [selectedElementType, setSelectedElementType] = useState(null);
    const { courseId, moduleId, submoduleId } = useParams();
    const [experiencePointReward, setExperiencePointReward] = useState(0);
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
        const fetchElements = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/module/${moduleId}/submodule/${submoduleId}/elements`);
                if (response.ok) {
                    const data = await response.json();
                    setElements(data);
                } else {
                    console.error('Failed to fetch elements');
                }
            } catch (error) {
                console.error('Error fetching elements:', error);
            }
        };
        fetchElements();
    }, [courseId, moduleId, submoduleId]);

    const handleVideoFileChange = (event) => {
        setVideoFile(event.target.files[0]);
    };

    const handleAddElement = async (type) => {
        let content = '';

        if (type === 'html') {
            content = htmlContent.trim();
        } else if (type === 'video') {
            if (videoFile) {
                const formData = new FormData();
                formData.append('video', videoFile);

                try {
                    const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/submodule/${submoduleId}/elements/video`, {
                        method: 'POST',
                        body: formData,
                    });

                    if (response.ok) {
                        const newElement = await response.json();
                        if (elements) {
                            setElements([...elements, newElement]);
                        } else {
                            setElements([newElement]);
                        }
                        setVideoFile(null);
                    } else {
                        console.error('Failed to upload video');
                    }
                } catch (error) {
                    console.error('Error uploading video:', error);
                }
                return;
            }
        } else if (type === 'quiz') {
            if (quizQuestion.trim() !== '' && quizAnswers.length > 0) {
                const quizData = {
                    submodule_id: parseInt(submoduleId),
                    question: quizQuestion,
                    question_type: quizType,
                    xp_reward: experiencePointReward,
                };

                try {
                    const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/quizzes`, {
                        method: 'POST',
                        headers: {
                            'Content-Type': 'application/json',
                        },
                        body: JSON.stringify(quizData),
                    });

                    if (response.ok) {
                        const newQuiz = await response.json();

                        const answerPromises = quizAnswers.map(async (answer, index) => {
                            const answerData = {
                                quiz_id: newQuiz.id,
                                answer_text: answer,
                                is_correct: quizType === 'single'
                                    ? answer === quizCorrectAnswer
                                    : quizCorrectAnswers.includes(answer),
                                order: index,
                            };

                            const answerResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/quiz/${newQuiz.id}/answers`, {
                                method: 'POST',
                                headers: {
                                    'Content-Type': 'application/json',
                                },
                                body: JSON.stringify(answerData),
                            });

                            if (!answerResponse.ok) {
                                throw new Error('Failed to create quiz answer');
                            }
                        });

                        await Promise.all(answerPromises);

                        setQuizQuestion('');
                        setQuizAnswers(['']);
                        setQuizCorrectAnswer('');
                        setQuizCorrectAnswers([]);
                        setExperiencePointReward(0);

                        // Create an element for the quiz
                        const elementData = {
                            submodule_id: parseInt(submoduleId),
                            type: quizType === 'single' ? 'quiz_single_choice' : 'quiz_multiple_choice',
                            quiz_id: newQuiz.id,
                            order: elements ? elements.length : 0,
                        };

                        const elementResponse = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/elements`, {
                            method: 'POST',
                            headers: {
                                'Content-Type': 'application/json',
                            },
                            body: JSON.stringify(elementData),
                        });

                        if (elementResponse.ok) {
                            const newElement = await elementResponse.json();
                            if (elements) {
                                setElements([...elements, newElement]);
                            } else {
                                setElements([newElement]);
                            }
                        } else {
                            console.error('Failed to create element for quiz');
                        }
                    } else {
                        console.error('Failed to create quiz');
                    }
                } catch (error) {
                    console.error('Error creating quiz:', error);
                }
            }
        }

        if (content !== '') {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/elements`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        submodule_id: parseInt(submoduleId),
                        type: type === 'quiz' ? (quizType === 'single' ? 'quiz_single_choice' : 'quiz_multiple_choice') : type,
                        content,
                        order: elements ? elements.length + 1 : 0,
                    }),
                });

                if (response.ok) {
                    const newElement = await response.json();
                    if (elements) {
                        setElements([...elements, newElement]);
                    } else {
                        setElements([newElement]);
                    }

                    if (type === 'html') {
                        setHtmlContent('');
                    } else if (type === 'quiz') {
                        setQuizQuestion('');
                        setQuizAnswers(['']);
                        setQuizCorrectAnswer('');
                        setQuizCorrectAnswers([]);
                        setExperiencePointReward(0);
                    }
                } else {
                    console.error('Failed to add element');
                }
            } catch (error) {
                console.error('Error adding element:', error);
            }
        }
    };

    const handleDeleteElement = async (id) => {
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/element/${id}`, {
                method: 'DELETE',
            });
            if (response.ok) {
                const updatedElements = elements.filter((element) => element.id !== id);
                setElements(updatedElements);
            } else {
                console.error('Failed to delete element');
            }
        } catch (error) {
            console.error('Error deleting element:', error);
        }
    };

    const onDragEnd = async (result) => {
        if (!result.destination) return;

        const newElements = Array.from(elements);
        const [reorderedElement] = newElements.splice(result.source.index, 1);
        newElements.splice(result.destination.index, 0, reorderedElement);

        try {
            const updatedOrder = newElements.reduce((acc, element, index) => {
                acc[element.id] = index;
                return acc;
            }, {});

            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/course/${courseId}/module/${moduleId}/submodule/${submoduleId}/elements/order`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ order: updatedOrder }),
            });
            if (response.ok) {
                setElements(newElements);
            } else {
                console.error('Failed to update element order');
            }
        } catch (error) {
            console.error('Error updating element order:', error);
        }
    };

    const handleRichTextEditorChange = (editorState) => {
        const contentState = editorState.getCurrentContent();
        const htmlContent = stateToHTML(contentState, {
            inlineStyles: {
                CODE: {
                    style: {
                        backgroundColor: "rgba(0, 0, 0, 0.05)",
                        fontFamily: '"Inconsolata", "Menlo", "Consolas", monospace',
                        fontSize: 16,
                        padding: 2,
                    }
                },
                HIGHLIGHT: {
                    style: {
                        backgroundColor: "#F7A5F7",
                    }
                },
                UPPERCASE: {
                    style: {
                        textTransform: "uppercase",
                    }
                },
                LOWERCASE: {
                    style: {
                        textTransform: "lowercase",
                    }
                },
                CODEBLOCK: {
                    style: {
                        fontFamily: '"fira-code", "monospace"',
                        fontSize: "inherit",
                        background: "#ffeff0",
                        fontStyle: "italic",
                        lineHeight: 1.5,
                        padding: "0.3rem 0.5rem",
                        borderRadius: " 0.2rem",
                    }
                },
                SUPERSCRIPT: {
                    style: {
                        verticalAlign: "super",
                        fontSize: "80%",
                    }
                },
                SUBSCRIPT: {
                    style: {
                        verticalAlign: "sub",
                        fontSize: "80%",
                    }
                },
            }, blockStyleFn: (contentBlock) => {
                const type = contentBlock.getType();
                switch (type) {
                    case "blockQuote":
                        return { attributes: { class: "superFancyBlockquote" } };
                    case "leftAlign":
                        return { attributes: { class: "leftAlign" } };
                    case "rightAlign":
                        return { attributes: { class: "rightAlign" } };
                    case "centerAlign":
                        return { attributes: { class: "centerAlign" } };
                    case "justifyAlign":
                        return { attributes: { class: "justifyAlign" } };
                    default:
                        break;
                }
            }
        });
        setHtmlContent(htmlContent);
    };

    if (!customSystemSettings) {
        return <div>Loading...</div>
    }

    return (
        <Layout user={user}>
            <Container maxWidth="lg" className="my-8">
                <Box mb={4}>
                    <Typography variant="h4" component="h1" gutterBottom>
                        Element Builder
                    </Typography>
                    <Box mb={2}>
                        <FormControl component="fieldset">
                            <FormLabel component="legend">Element type to add</FormLabel>
                            <RadioGroup
                                value={selectedElementType}
                                onChange={(e) => setSelectedElementType(e.target.value)}
                            >
                                <FormControlLabel value="html" control={<Radio />} label="HTML Content" />
                                <FormControlLabel value="video" control={<Radio />} label="Video" />
                                <FormControlLabel value="quiz" control={<Radio />} label="Quiz" />
                            </RadioGroup>
                        </FormControl>
                    </Box>
                    {selectedElementType === 'html' && (
                        <Box display="flex" alignItems="center" mb={2}>
                            <RichTextEditor onChange={handleRichTextEditorChange} />
                            <Button variant="contained" color="primary" onClick={() => handleAddElement('html')} style={{ marginLeft: '16px' }}>
                                Add
                            </Button>
                        </Box>
                    )}
                    {selectedElementType === 'video' && (
                        <Box display="flex" alignItems="center" mb={2}>
                            <input
                                type="file"
                                accept="video/*"
                                onChange={handleVideoFileChange}
                            />
                            <Button
                                variant="contained"
                                color="primary"
                                onClick={() => handleAddElement('video')}
                                style={{ marginLeft: '16px' }}
                                disabled={!videoFile}
                            >
                                Add
                            </Button>
                        </Box>
                    )}
                    {selectedElementType === 'quiz' && (
                        <Box mb={2}>
                            <TextField
                                label="Quiz Question"
                                value={quizQuestion}
                                onChange={(e) => setQuizQuestion(e.target.value)}
                                fullWidth
                                margin="normal"
                            />
                            <TextField
                                label="Quiz Answers (comma-separated)"
                                value={quizAnswers.join(',')}
                                onChange={(e) => setQuizAnswers(e.target.value.split(','))}
                                fullWidth
                                margin="normal"
                            />
                            <FormControl component="fieldset" margin="normal">
                                <FormLabel component="legend">Quiz Type</FormLabel>
                                <RadioGroup
                                    value={quizType}
                                    onChange={(e) => setQuizType(e.target.value)}
                                    row
                                >
                                    <FormControlLabel value="single" control={<Radio />} label="Single Choice" />
                                    <FormControlLabel value="multiple" control={<Radio />} label="Multiple Choice" />
                                </RadioGroup>
                            </FormControl>
                            {quizType === 'single' ? (
                                <TextField
                                    label="Correct Answer"
                                    value={quizCorrectAnswer}
                                    onChange={(e) => setQuizCorrectAnswer(e.target.value)}
                                    fullWidth
                                    margin="normal"
                                />
                            ) : (
                                <TextField
                                    label="Correct Answers (comma-separated)"
                                    value={quizCorrectAnswers.join(',')}
                                    onChange={(e) => setQuizCorrectAnswers(e.target.value.split(','))}
                                    fullWidth
                                    margin="normal"
                                />
                            )}

                            <TextField
                                label="Experience Point Reward"
                                type="number"
                                value={experiencePointReward}
                                onChange={(e) => setExperiencePointReward(parseInt(e.target.value))}
                                fullWidth
                                margin="normal"
                                inputProps={{
                                    min: 0,
                                    step: 1,
                                }}
                            />

                            <Button variant="contained" color="primary" onClick={() => handleAddElement('quiz')}>
                                Add Quiz
                            </Button>
                        </Box>
                    )}
                </Box>
                <DragDropContext onDragEnd={onDragEnd}>
                    <Droppable droppableId="elements">
                        {(provided) => (
                            <Box {...provided.droppableProps} ref={provided.innerRef}>
                                {elements && elements.map((element, index) => (
                                    <Draggable key={element.id} draggableId={element.id.toString()} index={index}>
                                        {(provided) => (
                                            <Box
                                                ref={provided.innerRef}
                                                {...provided.draggableProps}
                                                display="flex"
                                                alignItems="center"
                                                justifyContent="space-between"
                                                mb={2}
                                            >
                                                {element.type === 'html' && (
                                                    <div dangerouslySetInnerHTML={{ __html: element.content }} />
                                                )}
                                                {element.type === 'video' && (
                                                    <video
                                                        src={`/api/v1/element/video/${element.id}`}
                                                        type="video/mp4"
                                                        controls
                                                        style={{
                                                            width: '100%',
                                                            height: 'auto',
                                                            maxWidth: '100%',
                                                        }}
                                                        preload="auto"
                                                        playsInline
                                                    />
                                                )}
                                                {(element.type === 'quiz_single_choice' || element.type === 'quiz_multiple_choice') && (
                                                    <Box mb={2}>
                                                        <Quiz quizId={element.quiz_id} dryRun={true} />
                                                    </Box>
                                                )}
                                                <Box display="flex">
                                                    <IconButton
                                                        {...provided.dragHandleProps}
                                                        color="primary"
                                                        aria-label="drag"
                                                    >
                                                        <DragIndicatorIcon />
                                                    </IconButton>
                                                    <IconButton
                                                        edge="end"
                                                        color="secondary"
                                                        onClick={() => handleDeleteElement(element.id)}
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

export default ElementBuilder;
