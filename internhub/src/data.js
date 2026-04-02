export const COLORS = ['#FF3CAC','#7C3AED','#2B86C5','#059669','#F59E0B','#EF4444','#06B6D4','#8B5CF6'];
export const EMOJIS = ['🚀','💡','🔥','⚡','🌐','💻','🎯','🏆'];
export const getColor = n => COLORS[(n || '').charCodeAt(0) % COLORS.length];
export const getEmoji = n => EMOJIS[(n || '').charCodeAt(0) % EMOJIS.length];
