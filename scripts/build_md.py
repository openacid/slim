#!/usr/bin/env python
# coding: utf-8

import os
import jinja2
import subprocess

def render_j2(tmpl_path, tmpl_vars, output_path):

    def include_file(name):
        return jinja2.Markup(loader.get_source(env, name)[0])

    loader = jinja2.FileSystemLoader(searchpath='./')
    env = jinja2.Environment(loader=loader,
                                      undefined=jinja2.StrictUndefined)
    env.globals['include_file'] = include_file
    template = env.get_template(tmpl_path)

    txt = template.render(tmpl_vars)

    with open(output_path, 'w') as f:
        f.write(txt)


def command(cmd, *arguments, **options):

    close_fds = options.get('close_fds', True)
    cwd = options.get('cwd', None)
    shell = options.get('shell', False)
    env = options.get('env', None)
    if env is not None:
        env = dict(os.environ, **env)
    stdin = options.get('stdin', None)

    subproc = subprocess.Popen([cmd] + list(arguments),
                                 close_fds=close_fds,
                                 shell=shell,
                                 cwd=cwd,
                                 env=env,
                                 encoding='utf-8',
                                 stdin=subprocess.PIPE,
                                 stdout=subprocess.PIPE,
                                 stderr=subprocess.PIPE, )

    out, err = subproc.communicate(input=stdin)

    subproc.wait()

    if subproc.returncode != 0:
        raise Exception(subproc.returncode, out, err)

    return out


if __name__ == "__main__":
    pkg = command('go', 'list', '.')
    name = pkg.strip().split('/')[-1]
    tmpl_vars = {
            "name": name
    }

    fns = os.listdir('docs')
    for fn in fns:
        if fn.endswith('.md.j2'):
            render_j2('docs/'+fn, tmpl_vars, fn[:-3])
