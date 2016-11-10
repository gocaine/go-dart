package i18n

var frenchYaml = `
game:
  message:
    rank: Le joueur %d a terminé à la position %d
    disconnect: La cible %s a été déconnectée
    score:
      "Marqué : %d"
    winner:
      "Gagnant : %s"
    player:
      exists: Nom de joueur déjà pris
      notadded: Le joueur n'a pu être ajouté
      next: Joueur suivant
  error:
    onhold: Le jeu est en pause et ne peut pas recevoir de flechette
    notstarted: La partie n'est pas commencée ou terminée
    cantstart: La partie ne peut pas commencer
    sector:
      invalid: Ce secteur n'est pas valide
  countup:
    error:
      target: Target doit valoir 61 au minimum
    rules: Tous les joueurs commencent avec 0 points et doivent atteindre la \"Target\"
    options:
      target: le score à atteindre
  cricket:
    message:
      open:
        "Ouvert : %s"
      close:
        "Fermé : %s"
      hit:
        "Touché : %d x %s"
    error:
      incompatible: Les options CutThroat et NoScore sont incompatibles
    options:
      noscore: Si cette option vaut true, aucun point n'est marqué, la gagnant est le premier joueur à ouvrir tous les secteurs
      cutthroat: Si cette option vaut true, à partir de la 4ème touche dans un secteur, les points sont attribués à tous les joueurs qui n'ont pas encore ouvert ce secteur. A la fin, le gagnant est le joueur qui le premier a ouvert tous les secteurs et possède le moins de points
    rules: Le principe du jeu est d'ouvrir (ou fermer) tous les secteurs. Les secteurs sont 15, 16, 17, 18, 19, 20 et la bulle. Pour ouvrir un secteur, un joueur doit le toucher 3 fois (un Triple compte pour 3 touches, un Double pour 2). Quand un secteur est ouvert pour un joueur, il peut marquer dedans (les points sont la valeur réelle de la touche). Quand tous les joueurs ont ouvert un secteur donné, il est alors fermé, et plus aucun point ne peut y être marqué. Le gagnant est le premier joueur a avoir simultanément ouvert tous les secteurs et possède le plus gros score
  highest:
    display: "HighScore sur %d volées"
    error:
      rounds: Le nombre de Rounds doit être 1 au minimum
    options:
      rounds: Le nombre de volée pour chaque joueur
    rules: Tous les joueurs lancent le même nombre de fléchettes (3 par volée) puis le joueur avec le score le plus haut gagne
  x01:
    message:
      doubleout: Vous devez finir par un double
      overscore: Vous avez dépassé !
    error:
      score: Le Score doit être 61 au minimum
    options:
      score: Le score de départ
      doubleout: Si cette option vaut true, les joueurs doivent finir par un double
    rules:
      Tous les joueurs démarrent avec les mêmes points (301 / 501 / ...) et doivent atteindre 0. Chaque flechette réduit ainsi le score. Si un joueur marque plus de points que son score restant alors il "bust" et il retourne a son score de début de volée
`

func init() {
	yamlFiles["fra"] = frenchYaml
}
