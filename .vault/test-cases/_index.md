---
type: index
title: Index des test cases
tags: [test-case, index]
---

# Index des test cases

- [[TC-00-initialisation]] — Le vault est complet et navigable, la chaîne de test fonctionne
- [[TC-01-api-rest-auth]] — Le pont rejette les requêtes non authentifiées, cross-origin ou non JSON
- [[TC-02-media-path-allowlist]] — Seuls les fichiers du dossier d'envoi peuvent être lus
- [[TC-03-download-path-traversal]] — Aucun nom malveillant ne sort du dossier du chat
- [[TC-04-allowed-chats]] — Les conversations hors liste blanche sont ignorées et masquées
- [[TC-05-read-only-mode]] — Les outils d'envoi disparaissent en lecture seule
- [[TC-06-store-permissions]] — store/ en 0700, fichiers en 0600
- [[TC-07-bump-whatsmeow]] — Le pont compile et fonctionne avec whatsmeow à jour
- [[TC-08-evenements-inbox-webhook]] — Un message autorisé produit média + événement JSON + appel webhook
- [[TC-09-api-envoi-health]] — Envoi idempotent, audité, réservé au jeton d'envoi ; health OK
- [[TC-10-evenements-edit-revoke-reaction]] — Edit / revoke / reaction produisent le bon type d'événement
- [[TC-11-integration-rapport]] — main du fork vert, PR ouvertes
- [[TC-12-code-appairage]] — Les codes d'appairage sont livrés, limités en fréquence, puis confirmés
- [[TC-13-api-history]] — L'historique est paginé, filtré par la liste blanche et au format v1
- [[TC-14-store-dir]] — Toutes les données vont dans WHATSAPP_STORE_DIR
- [[TC-15-appairage-robuste]] — L'appairage livre toujours ses événements et s'arrête proprement
- [[TC-16-confidentialite-journaux-stockage]] — Aucun texte journalisé, liste blanche stricte, purge effective
- [[TC-17-api-groups]] — La liste des groupes ne contient que les métadonnées autorisées
- [[TC-18-reaction-reponse-citee]] — Réactions et réponses citées visent le bon message, sous les mêmes garde-fous que l'envoi
- [[TC-19-legende-medias]] — La légende d'un média est stockée comme contenu du message
- [[TC-20-numero-expediteur]] — Le numéro de l'expéditeur est stocké même quand WhatsApp n'envoie que le LID
- [[TC-21-telechargement-media]] — Le téléchargement d'un média réussit avec le chemin signé, la copie de l'inbox ou un renouvellement
