#!/usr/bin/env python

try:
    from setuptools import setup
except ImportError:
    from distutils.core import setup

VERSION = '2.4.1'

setup(name='tornado-redis',
      version=VERSION,
      description='Asynchronous Redis client for the Tornado Web Server.',
      author='Vlad Glushchuk',
      author_email='vgluschuk@gmail.com',
      license="http://www.apache.org/licenses/LICENSE-2.0",
      url='http://github.com/leporo/tornado-redis',
      keywords=['Redis', 'Tornado'],
      packages=['tornadoredis'], )
