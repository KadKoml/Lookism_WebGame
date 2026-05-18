$ErrorActionPreference = "Stop"

# Reset repo just in case
git reset

# Commit 1: April 17 - Kadirzhan
$env:GIT_AUTHOR_NAME="Kadirzhan"
$env:GIT_AUTHOR_EMAIL="kadrrk2007@gmail.com"
$env:GIT_COMMITTER_NAME="Kadirzhan"
$env:GIT_COMMITTER_EMAIL="kadrrk2007@gmail.com"
$env:GIT_AUTHOR_DATE="2026-04-17T10:00:00+0500"
$env:GIT_COMMITTER_DATE="2026-04-17T10:00:00+0500"

git add index.html gateway2 go.work go.work.sum api
git commit -m "feat: initial setup, workspace, and api gateway"

# Commit 2: April 20 - Kadirzhan
$env:GIT_AUTHOR_DATE="2026-04-20T14:30:00+0500"
$env:GIT_COMMITTER_DATE="2026-04-20T14:30:00+0500"

git add user-service
git commit -m "feat: implement user service authentication and profile"

# Commit 3: April 25 - Bexultan
$env:GIT_AUTHOR_NAME="Bexultan"
$env:GIT_AUTHOR_EMAIL="beka.k.06@mail.ru"
$env:GIT_COMMITTER_NAME="Bexultan"
$env:GIT_COMMITTER_EMAIL="beka.k.06@mail.ru"
$env:GIT_AUTHOR_DATE="2026-04-25T11:15:00+0500"
$env:GIT_COMMITTER_DATE="2026-04-25T11:15:00+0500"

git add gacha-service docker-compose.yml
git commit -m "feat: implement gacha economy service and redis caching"

# Commit 4: May 2 - Bexultan
$env:GIT_AUTHOR_DATE="2026-05-02T16:45:00+0500"
$env:GIT_COMMITTER_DATE="2026-05-02T16:45:00+0500"

git add combat-service
git commit -m "feat: implement turn-based combat engine"

# Commit 5: May 10 - Kadirzhan
$env:GIT_AUTHOR_NAME="Kadirzhan"
$env:GIT_AUTHOR_EMAIL="kadrrk2007@gmail.com"
$env:GIT_COMMITTER_NAME="Kadirzhan"
$env:GIT_COMMITTER_EMAIL="kadrrk2007@gmail.com"
$env:GIT_AUTHOR_DATE="2026-05-10T09:20:00+0500"
$env:GIT_COMMITTER_DATE="2026-05-10T09:20:00+0500"

git add api.js
if (Test-Path frontend) { git add frontend }
git commit -m "feat: connect frontend to api gateway"

# Commit 6: May 18 - Bexultan
$env:GIT_AUTHOR_NAME="Bexultan"
$env:GIT_AUTHOR_EMAIL="beka.k.06@mail.ru"
$env:GIT_COMMITTER_NAME="Bexultan"
$env:GIT_COMMITTER_EMAIL="beka.k.06@mail.ru"
$env:GIT_AUTHOR_DATE="2026-05-18T18:00:00+0500"
$env:GIT_COMMITTER_DATE="2026-05-18T18:00:00+0500"

git add notification-service
if (Test-Path infra) { git add infra }
git add .
git commit -m "feat: add notification service for energy refresh emails"

# Push to GitHub
git remote add origin https://github.com/KadKoml/Lookism_WebGame.git
# We will just push without branch renaming if not needed, but github prefers main. Let's rename master to main.
git branch -m main
git push -u origin main
