package comics

var stopwords map[string]bool

func init() {
	stopwords = make(map[string]bool)

	stopWordsList := []string{
		"a", "about", "above", "after", "again", "against", "all", "am", "an",
		"and", "any", "are", "as", "at", "be", "because", "been", "before",
		"being", "below", "between", "both", "but", "by", "can", "did", "do",
		"does", "doing", "don", "down", "during", "each", "few", "for", "from",
		"further", "had", "has", "have", "having", "he", "her", "here", "hers",
		"herself", "him", "himself", "his", "how", "i", "if", "in", "into", "is",
		"it", "its", "itself", "just", "me", "more", "most", "my", "myself",
		"no", "nor", "not", "now", "of", "off", "on", "once", "only", "or",
		"other", "our", "ours", "ourselves", "out", "over", "own", "s", "same",
		"she", "should", "so", "some", "such", "t", "than", "that", "the", "their",
		"theirs", "them", "themselves", "then", "there", "these", "they",
		"this", "those", "through", "to", "too", "under", "until", "up",
		"very", "was", "we", "were", "what", "when", "where", "which", "while",
		"who", "whom", "why", "will", "with", "you", "your", "yours", "yourself",
		"yourselves", "i'v", "i'm", "i'd", "i'll", "i've", "you're", "you'd", "you'll",
		"you've", "he's", "he'd", "he'll", "she's", "she'd", "she'll", "it's", "it'll",
		"we're", "we'd", "we'll", "we've", "they're", "they'd", "they'll", "they've",
		"there's", "there'll", "there'd", "isn't", "aren't", "don't", "doesn't", "wasn't",
		"weren't", "didn't", "haven't", "hasn't", "won't", "hadn't", "can't", "couldn't",
		"mustn't", "mightn't", "needn't", "shouldn't", "oughtn't", "wouldn't",
		"what's", "how's", "where's", "we'r",
	}

	for _, word := range stopWordsList {
		stopwords[word] = true
	}
}

func IsStopWord(word string) bool {
	return stopwords[word]
}
