import React, { useState, useEffect } from 'react';
import { FormControl, FormLabel, FormGroup, FormControlLabel, Checkbox, Button } from '@mui/material';

const MultipleChoiceQuiz = ({ elementId }) => {
    const [quiz, setQuiz] = useState(null);
    const [selectedAnswers, setSelectedAnswers] = useState([]);
    const [result, setResult] = useState(null);

    useEffect(() => {
        const fetchQuiz = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/element/${elementId}`);
                const data = await response.json();
                const responseQuiz = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/quiz/${data.quiz_id}`, {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    }
                });
                const quizData = await responseQuiz.json();
                setQuiz(quizData);
            } catch (error) {
                console.error('Error fetching quiz:', error);
            }
        };
        fetchQuiz();
    }, [elementId]);

    const handleChange = (option) => {
        const index = selectedAnswers.indexOf(option);
        if (index === -1) {
            setSelectedAnswers([...selectedAnswers, option]);
        } else {
            setSelectedAnswers(selectedAnswers.filter((answer) => answer !== option));
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/checkQuiz/${quiz.id}`, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                },
                body: JSON.stringify(selectedAnswers),
            });
            const data = await response.json();
            setResult(data.correct);
        } catch (error) {
            console.error('Error checking quiz answers:', error);
        }
    };

    if (!quiz) {
        return <div>Loading...</div>;
    }

    return (
        <form onSubmit={handleSubmit}>
            <FormControl component="fieldset">
                <FormLabel component="legend">{quiz.question}</FormLabel>
                <FormGroup>
                    {quiz.quiz_answers.map((answer) => (
                        <FormControlLabel
                            key={answer.id}
                            control={
                                <Checkbox
                                    checked={selectedAnswers.includes(answer.answer_text)}
                                    onChange={() => handleChange(answer.answer_text)}
                                    disabled={quiz.answered_by_user || result}
                                />
                            }
                            label={answer.answer_text}
                        />
                    ))}
                </FormGroup>
                <Button type="submit" variant="outlined" sx={{ mt: 1, mr: 1 }} disabled={quiz.answered_by_user || result}>
                    Submit
                </Button>
            </FormControl>
            {(result !== null || quiz.answered_by_user) && (
                <div>
                    {result || quiz.answered_by_user ? 'You got it!' : 'Sorry, wrong answer!'}
                </div>
            )}
        </form>
    );
};

export default MultipleChoiceQuiz;