#!/usr/bin/env python3
# In The Name Of God
# ========================================
# [] File Name : main.py
#
# [] Creation Date : 04-02-2018
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import click
import runpy
import json

from actuator import Actuator


@click.command()
@click.argument('target', type=click.Path())
@click.argument('state')
@click.argument('action')
def run(target,state,action):
    '''
    run given target in provided environment
    '''
    try:
        g = runpy.run_path(target, run_name='uactuator')
        for value in g.values():
            if isinstance(value, type) and issubclass(value, Actuator) and \
                    value.__module__ == 'uactuator':
                actuator = value
    except Exception as e:
        print('Target Error: ', e)
        return

    d = actuator().value(state,action)
    print(json.dumps(d))

def main():
    run()
