#!/bin/sh
#  @brief tiny tool used to manger project 
#  @file repos.sh
#  @author xiyanxiyan10
#  @CreatedTime 2019/08/09
#  @LastChanged 2019/08/09

case "$1" in
    "record" )
        git commit -a
        ;;
    "pull" )
        git pull origin master
        ;;
    "push" )
        git push 
        ;;
    "changelog" )
        git log --oneline
        ;;
    "checkout" )
        git checkout -- .
        ;;
    "revert" )
        git revert ${2}
        ;;
    "add_remote" )
        git remote add ${2} git@github.com:${2}/${3}.git
        ;;
    "pull_remote" )
        git remote add ${2} git@github.com:${2}/${3}.git || git fetch ${2} && git merge ${2}/master
        ;;
    "fetch" )
        git fetch ${2}
        ;;
    "merge" )
        git merge ${2}/master
        ;;
    "show_remote"  )
        git remote -v
        ;;
    "tag" )
        git tag -a ${2}
        ;;
    "pushtag" )
        git push git@github.com:${2}/${3}.git --tag
        ;;
    "clean" )
        rm -rf site
        find ./ -name '*.so'  |xargs rm 
        find ./ -name '*.o'   |xargs rm 
        find ./ -name '*.a'   |xargs rm 
        find ./ -name '*.d'   |xargs rm 
        find ./ -name '*.pyc' |xargs rm 
        ;;
    "review")
        cscope -bqR
        ctags -R *
        ;;
    "pyformat")
        yapf -p --style='{based_on_style: chromium, indent_width: 4}' -i -r ./
        ;;
    "replace")
        ${2} | xargs perl -pi -e \'${3}\'
        ;;
    "apidocs")
        swagger_py_codegen -s api.yaml site -p site -tlp=tornado --ui --spec  
        ;;
    * ) 
        echo "Help"
        echo "./repos.sh apidocs                                => build swaggar api files"
        echo "./repos.sh pyformat                               => format project python files"
        echo "./repos.sh record                                 => record patch"
        echo "./repos.sh pull                                   => pull patch"
        echo "./repos.sh push                                   => push patch"
        echo "./repos.sh changelog                              => show changelog"
        echo "./repos.sh checkout                               => revert change code"
        echo "./repos.sh revert                                 => revert patch by id"
        echo "./repos.sh add_remote                             => add new remote"
        echo "./repos.sh fetch                                  => fetch from remote"
        echo "./repos.sh merge                                  => merge from remote"
        echo "./repos.sh pull_remote                            => pull from remote"
        echo "./repos.sh show_remote                            => show remote"
        echo "./repos.sh tag                                    => tag version"
        echo "./repos.sh pushtag                                => push tag"
        echo "./repos.sh clean                                  => clean project"      
        echo "./repos.sh review                                 => review project"
        echo "./repos.sh [cmd] [replace pattern]                => change some pattern matched in matched file"
        ;;
    esac
