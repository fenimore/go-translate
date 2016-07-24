# go-translate
A command line tool for finding English to French and French to English translations.

## Who says?
I use the Glosbe translation API and the wordreference website to get conjugation/word parellels.

## Use
I put the binary inside `$GOPATH/bin`, export this path (or create an alias). So either run `def e word` for en-fr or `def f word` for fr-en; or create an alias for enfr or fren. The second argument passed to def will determine the origin language.

## Example 

Run 

    def f blague
	
And it outputs (what is in practice a color output) the translated words, the possible conjugations and up to 30 sentences in French and English which use the given word.

	joke, trick, pouch, mess about, muck around, banter, josh, cheap joke, corny joke, joking aside, joking apart, being serious, dirty joke, schoolboy prank, below-the-belt joke, I'm joking!, I'm kidding!, Just kidding!, bad joke, sick joke, Seriously?! Really?!, no more joking around, no more kidding around, no more joking about, no more kidding about, seriously, seriously though, Really!, Honestly!, No kidding! Joking aside!, 

	FR-EN:     blague 
	Translate: joke, trick 
	Du verbe blaguer -> blague est:
	1 re personne du singulier du présent de l'indicatif 3 e personne du singulier du présent de l'indicatif 1 re personne du singulier du présent du subjonctif 3 e personne du singulier du présent du subjonctif 2 e personne du singulier du présent de l'impératif blagué est:
	un participe passé 
	From: 
	Je blague pas
	To:
	I' m dead serious
	More 1/30? [y] 


