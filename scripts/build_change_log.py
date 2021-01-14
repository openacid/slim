#!/usr/bin/env python
# coding: utf-8

import sys
import subprocess
import semantic_version
import yaml
from collections import defaultdict

with open('scripts/change-types.yaml', 'r') as f:
    typs = yaml.load(f.read())

typs = {x:x for x in typs}

# categories has another mapping to fix typo in commit message
categories = {
        'api-change:':   typs['api-change'],
        'new-feature:':  typs['new-feature'],
        'internal:':     typs['internal'],
        'doc:':          typs['doc'],
        'refactor:':     typs['refactor'],
        'fixbug:':       typs['fixbug'],
        'fixdoc:':       typs['fixdoc'],

        # fix typo
        'api-changes:':  typs['api-change'],
        'new-features:': typs['new-feature'],
        'docs:':         typs['doc'],
        'fix:':          typs['fixbug'],

        'test:':         typs['test'],
}

to_display = {
        'doc': False,
        'refactor': False,
        'internal': False,
        'test': False,
}

def cmd(cmds):
    subproc = subprocess.Popen(cmds,
                               encoding='utf-8',
                                  stdout=subprocess.PIPE,
                                  stderr=subprocess.PIPE, )
    out, err = subproc.communicate()
    subproc.wait()

    code = subproc.returncode
    if code != 0:
        raise OSError(out + "\n" + err)

    return out

def list_tags():
    out = cmd(["git", "tag", "-l"])
    tags = out.splitlines()
    tags[0].lstrip('v')
    tags = [semantic_version.Version(t.lstrip('v'))
            for t in tags
            if t != '' ]
    return tags


def changes(frm, to):
    # subject, author time, author name, email.
    out = cmd(["git", "log", '--format=%s ||| %ai ||| %an ||| %ae', '--reverse', frm + '..' + to])
    lines = out.splitlines()
    lines = [x for x in lines if x != '']
    rst = []
    for line in lines:
        elts = line.split(" ||| ")
        item = {
                'subject': elts[0],
                # 2019-04-18 13:36:42 +0800
                'time': elts[1].split()[0],
                'author': elts[2],
                'email': elts[3]
        }
        rst.append(item)

    return rst

def norm_changes(changes):
    rst = {}
    for ch in changes:
        sub = ch['subject']
        cate, mod, cont = sub.split(' ', 2)
        catetitle = categories[cate]
        mod = mod.rstrip(':')

        if catetitle not in rst:
            rst[catetitle] = {}

        c = rst[catetitle]
        if mod not in c:
            c[mod] = []

        l = c[mod]
        desc = {
                "content": cont.replace(':', ''),
                "time": ch['time'],
                "author": ch['author'],
                'email': ch['email'],
        }
        l.append('{content}; by {author}; {time}'.format(**desc))

    return rst

def build_ver_changelog(newver):
    tags = list_tags()
    tags.sort()

    newver = newver.lstrip('v')
    newver = semantic_version.Version(newver)
    tags = [t for t in tags if t < newver]
    latest = tags[-1]

    chs = changes('v' + str(latest), 'HEAD')
    chs = norm_changes(chs)
    chs = {k:v for k, v in chs.items() if to_display.get(k, True) }

    changelog = yaml.dump(chs, default_flow_style=False)

    with open('docs/change-log/v{newver}.yaml'.format(newver=newver), 'w') as f:
        f.write(changelog)

def build_changelog():

    out = cmd(["ls", "docs/change-log"])
    vers = out.splitlines()
    # remove "yaml"
    vers = [x.rsplit('.', 1)[0] for x in vers if x != '']
    vers.sort(key=lambda x: semantic_version.Version(x.lstrip('v')))

    with open('docs/change-log.yaml', 'w') as f:
        for v in reversed(vers):
            f.write(v + ':\n')
            with open('docs/change-log/{v}.yaml'.format(v=v), 'r') as vf:
                cont = vf.read()

            cont = cont.splitlines()
            cont = ['  ' + x for x in cont]
            cont = '\n'.join(cont)

            f.write(cont + '\n')

if __name__ == "__main__":
    # Usage: to build change log from git log
    # ./scripts/build_change_log.py v0.5.10
    newver = sys.argv[1]
    build_ver_changelog(newver)
    build_changelog()

