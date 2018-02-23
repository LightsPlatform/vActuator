# In The Name Of God
# ========================================
# [] File Name : actuator/actuator.py
#
# [] Creation Date : 04-02-2018
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc


class Actuator(metaclass=abc.ABCMeta):
    @abc.abstractmethod
    def value(self,state):
        pass
