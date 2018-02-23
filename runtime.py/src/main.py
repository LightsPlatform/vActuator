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
def run(target,state):
    '''
    run given target in provided environment
    '''
    try:
        g = runpy.run_path(target, run_name='actuator')
        for value in g.values():
            if isinstance(value, type) and issubclass(value, actuator) and \
                    value.__module__ == 'actuator':
                actuator = value
    except Exception as e:
        print('Target Error: ', e)
        return

    d = actuator().value(state)
    print(json.dumps(d))

def main():
    run()
