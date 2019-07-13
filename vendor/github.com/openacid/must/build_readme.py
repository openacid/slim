#!/usr/bin/env python
# coding: utf-8


import jinja2

def rr(tmpl_path, tmpl_vars, output_path):
    template_loader = jinja2.FileSystemLoader(searchpath='./')
    template_env = jinja2.Environment(loader=template_loader,
                                      undefined=jinja2.StrictUndefined)
    template = template_env.get_template(tmpl_path)

    txt = template.render(tmpl_vars)

    with open(output_path, 'w') as f:
        f.write(txt.encode('utf-8'))

if __name__ == "__main__":
    rr('README.md.j2', {}, 'README.md')
