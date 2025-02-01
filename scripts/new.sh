cd "$(dirname "$0")" || exit
cd ..
touch changelog.md

echo -e "\n$(cat changelog.md)\n\n-----\n\n**[ $(date +"%Y-%m-%d %H:%M") ]**: $1" > changelog.md

git add .
git commit -m "$1"