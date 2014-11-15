import os
from setuptools import setup, find_packages
here = os.path.abspath(os.path.dirname(__file__))

setup(name='aguin',
      version='0.1',
      author='Van Nhu Nguyen',
      author_email='',
      url='https://github.com/vanng822/aguin',
      packages=find_packages(),
      test_suite="tests",
      install_requires=['requests', 'M2Crypto']
)