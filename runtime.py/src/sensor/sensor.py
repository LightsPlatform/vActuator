# In The Name Of God
# ========================================
# [] File Name : sensor/sensor.py
#
# [] Creation Date : 04-02-2018
#
# [] Created By : Parham Alvani <parham.alvani@gmail.com>
# =======================================
import abc


class Sensor(metaclass=abc.ABCMeta):
    @abc.abstractmethod
    def value(self):
        pass
