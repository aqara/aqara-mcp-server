git remote add origin https://github.com/aqara/aqara-mcp-server.git
git config --local user.name "Eason"
git config --local user.email "kingbigcountry@proton.me"

# git config --local http.proxy socket5://127.0.0.1:33210
git config --local --unset http.proxy
# git config --local https.proxy socket5://127.0.0.1:33210
git config --local --unset https.proxy

git checkout --orphan temp
git add .
git commit -m "first commit"
git branch -D main
git branch -m main
git push -f origin main

# 立即清除所有分支和标签的引用日志
git reflog expire --expire=now --all
# 立即清理所有未引用的对象，并以更积极的方式优化仓库
git gc --prune=now --aggressive

