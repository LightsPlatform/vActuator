# In The Name Of God
# ========================================
# [] File Name : setup.py
#
# [] Creation Date : 16-11-2017
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
'''
setup module for runtime.py
'''

# To use a consistent encoding
from os import path
# Always prefer setuptools over distutils
from setuptools import setup
# parse requirements.txt
from pip.req import parse_requirements

HERE = path.abspath(path.dirname(__file__))

# parse_requirements() returns generator of pip.req.InstallRequirement objects
INSTALL_REQS = parse_requirements(path.join(HERE, 'requirements.txt'),
                                  session='hack')

# reqs is a list of requirement
# e.g. ['django==1.5.1', 'mezzanine==1.4.6']
REQS = [str(ir.req) for ir in INSTALL_REQS]


setup(
    name='runtime.py',

    # Versions should comply with PEP440.
    # For a discussion on single-sourcing
    # the version across setup.py and the project code, see
    # https://packaging.python.org/en/latest/single_source_version.html
    version='0.1.0',


    # Author details
    author='Parham Alvani',
    author_email='parham.alvani@gmail.com',

    package_dir={'': 'src'},
    packages=['sensor'],
    py_modules=['main'],

    # List run-time dependencies here.  These will be installed by pip when
    # your project is installed.
    # For an analysis of "install_requires" vs pip's
    # requirements files see:
    # https://packaging.python.org/en/latest/requirements.html
    install_requires=REQS,

    # To provide executable scripts, use entry points in preference to the
    # "scripts" keyword. Entry points provide cross-platform
    # support and allow
    # pip to create the appropriate form of executable
    # for the target platform.
    entry_points={
        'console_scripts': [
            'runtime.py=main:main',
        ],
    },

)
