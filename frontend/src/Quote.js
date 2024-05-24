
const quotes = [
    "The only way to do great work is to love what you do.",
    "Believe you can and you're halfway there.",
    "Don't watch the clock; do what it does. Keep going.",
    "The harder you work for something, the greater you'll feel when you achieve it.",
    "Success is not the key to happiness. Happiness is the key to success. If you love what you are doing, you will be successful."
];

function getQuoteOfTheDay() {
    const today = new Date();
    const dayOfYear = Math.floor((Date.UTC(today.getFullYear(), today.getMonth(), today.getDate()) - Date.UTC(today.getFullYear(), 0, 0)) / (24 * 60 * 60 * 1000));
    return quotes[dayOfYear % quotes.length];
}

export default getQuoteOfTheDay;