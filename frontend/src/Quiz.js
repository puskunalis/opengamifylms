import React, { useState, useEffect } from 'react';
import { FormControl, FormLabel, RadioGroup, Radio, FormControlLabel, Checkbox, FormGroup, Button } from '@mui/material';

const Quiz = ({ quizId, dryRun, onQuizComplete }) => {
    const [quiz, setQuiz] = useState(null);
    const [selectedAnswers, setSelectedAnswers] = useState([]);
    const [result, setResult] = useState(null);

    useEffect(() => {
        if (quiz && quiz.answered_by_user && onQuizComplete) {
            onQuizComplete(quiz.id, true);
        }
    }, [quiz, onQuizComplete]);

    useEffect(() => {
        const fetchQuiz = async () => {
            try {
                const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/quiz/${quizId}`, {
                    headers: {
                        'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    }
                });
                const quizData = await response.json();
                setQuiz(quizData);
            } catch (error) {
                console.error('Error fetching quiz:', error);
            }
        };
        fetchQuiz();
    }, [quizId]);

    const handleChange = (answer) => {
        if (quiz.question_type === 'single') {
            setSelectedAnswers([answer]);
        } else {
            const index = selectedAnswers.indexOf(answer);
            if (index === -1) {
                setSelectedAnswers([...selectedAnswers, answer]);
            } else {
                setSelectedAnswers(selectedAnswers.filter((item) => item !== answer));
            }
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await fetch(`${process.env.REACT_APP_CUSTOM_BACKEND || ''}/api/v1/checkQuiz/${quiz.id}?dryRun=${dryRun}`, {
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
                {quiz.question_type === 'single' ? (
                    <RadioGroup value={selectedAnswers[0] || ''}>
                        {quiz.quiz_answers.map((answer) => (
                            <FormControlLabel
                                key={answer.id}
                                value={answer.answer_text}
                                control={<Radio onChange={() => handleChange(answer.answer_text)} disabled={quiz.answered_by_user || result} />}
                                label={answer.answer_text}
                            />
                        ))}
                    </RadioGroup>
                ) : (
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
                )}
                <Button type="submit" variant="outlined" sx={{ mt: 1, mr: 1 }} disabled={quiz.answered_by_user || result}>
                    Submit
                </Button>
            </FormControl>
            {(result !== null || quiz.answered_by_user) && (
                <div>
                    {result || quiz.answered_by_user ? `You got it! ${quiz.xp_reward > 0 ? `+${quiz.xp_reward} XP` : ''}` : 'Sorry, wrong answer!'}
                </div>
            )}
        </form>
    );
};

export default Quiz;