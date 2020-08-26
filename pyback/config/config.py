import yaml
import os

class ConfigService(object):
    """yaml config loader"""

    def __init__(self):
        """Init config Instance."""
        self.data = None

    def load(self, yaml_file):
        """Load configfile into parser"""
        file = open(yaml_file, 'r')
        file_data = file.read()
        file.close()

        data = yaml.load(file_data, Loader=yaml.FullLoader)
        #print(data)
        self.data = data
        return data

    def get(self, key):
        """get val by key, config only with map in root path"""
        if self.data == None:
            return None
        else:
            return self.data[key]


configService = ConfigService()


def GetConfig():
    return configService


if __name__ == '__main__':
    config = GetConfig()
    config.load("/tmp/config.yaml")
    print(config.get("test"))
    config = GetConfig()
    print(config.get("test"))
