import logging
import os
import sys
import traceback
try:
    from fdfs_client.client import *
    from fdfs_client.exceptions import *
except ImportError:
    import_path = os.path.abspath('../')
    sys.path.append(import_path)
    from fdfs_client.client import *
    from fdfs_client.exceptions import *

import setting

class DfsUtils():
    def __init__(self, conf_name=None):
        base_file_path =  os.path.dirname(__file__)

        fastdfs_client_conf_path = os.path.join(base_file_path, "../fastdfs_client.conf")
        self.mongo_host = setting.MONGO_HOST
        self.mongo_port = setting.MONGO_PORT

        if conf_name:
            self.client = Fdfs_client(conf_name)
        else:
            self.client = Fdfs_client(fastdfs_client_conf_path)

    """ check_file_exists: return fileid """
    def check_file_exists(self, file_md5):
        return None
        """
        from pymongo import Connection
        connection = Connection(self.mongo_host, self.mongo_port)
        db = connection.files
        files_data = db.files.find_one({"md5": file_md5})
        if files_data:
            return files_data["fileid"]
        else:
            return None
        """


    """ add file check data """
    def add_file_check_data(self, file_md5, fileid):
        return
        """
        from pymongo import Connection
        connection = Connection(self.mongo_host, self.mongo_port)
        db = connection.files
        db.files.insert({"md5": file_md5, "fileid": fileid})
        """

    """ delete file check data """
    def del_file_check_data(self, file_md5):
        return
        """
        from pymongo import Connection
        connection = Connection(self.mongo_host, self.mongo_port)
        db = connection.files
        db.files.remove({"md5": file_md5})
        """

    """
    delele_by_filename
    """
    def delete_by_filename(self, file_md5, remote_fileid):
        return
        self.del_file_check_data(file_md5)
        try:
            self.client.delete_file(str(remote_fileid))
        except Exception as inst:
            logging.error(inst.args)

    def delete_file(self, remote_fileid):
        try:
            return self.client.delete_file(str(remote_fileid))
        except Exception as inst:
            logging.debug(inst.args)
            return None

    """
    upload_by_filename:
    @meta_dict usage --
        meta_dict = {
            'ext_name' : 'jpg',
            'file_size' : '128KB'
        }
    """
    def upload_by_filename(self, file_md5, local_filename, meta_dict=None):
        fileid = self.check_file_exists(file_md5)
        if fileid == None:
            ret_dict = self.client.upload_by_filename(local_filename, meta_dict)
            if ret_dict:
                """ for test """
                for key in ret_dict:
                    logging.debug('[+] %s : %s' % (key, ret_dict[key]))
                """ add file check data """
                self.add_file_check_data(file_md5, ret_dict["Remote file_id"])
                return ret_dict["Remote file_id"]
            else:
                return None
        else:
            """ for test """
            logging.debug("file exists!")
            return fileid

    """
    upload_by_buffer:
    @file_buffer usage --
        with open(local_filename, 'rb') as f:
            file_buffer = f.read()
    @ext_name --
        like jpg png
    @meta_buffer usage --
        meta_buffer = {
            'ext_name' : 'jpg',
            'width' : '256',
            'height' : '256'
        }
    """
    def upload_by_buffer(self, file_md5, file_buffer, ext_name=None, meta_buffer=None):
        fileid = self.check_file_exists(file_md5)
        if fileid == None:
            ret_dict = self.client.upload_by_buffer(file_buffer, ext_name, meta_buffer)
            if ret_dict:
                """ for test """
                for key in ret_dict:
                    logging.debug('[+] %s : %s' % (key, ret_dict[key]))
                """ add file check data """
                self.add_file_check_data(file_md5, ret_dict["Remote file_id"])
                return ret_dict["Remote file_id"]
            else:
                return None
        else:
            """ for test """
            logging.debug("file exists!")
            return fileid

    """
    download_to_file:
    @local_filename string -- local filepath
    @remote_fileid string -- file identifier
    """
    def download_to_file(self, local_filename, remote_fileid):
        ret_dict = self.client.download_to_file(local_filename, remote_fileid)
        """ for test """
        for key in ret_dict:
            logging.debug('[+] %s : %s' % (key, ret_dict[key]))

    """
    download_to_buffer:
    @remote_fileid string -- file identifier
    @return buffer -- file buffer
    """
    def download_to_buffer(self, remote_fileid):
        ret_dict = self.client.download_to_buffer(remote_fileid)
        return ret_dict['Content']

    def download_to_buffer_all(self, remote_fileid):
        ret_dict = self.client.download_to_buffer(remote_fileid)
        return ret_dict

    """
    upload_slave_by_filename:
    """
    def upload_slave_by_filename(self, file_md5, local_filename, remote_file_id, prefix_name, \
                                 meta_dict=None, check_status=True):
        file_md5 = file_md5 + prefix_name
        fileid = None
        fileid = self.check_file_exists(file_md5)
        if check_status == False:
            self.delete_by_filename(file_md5, fileid)
            fileid = None
        if fileid is None:
            try:
                ret_dict = self.client.upload_slave_by_file(local_filename, remote_file_id, \
                                                            prefix_name, meta_dict)
                if ret_dict:
                    """ for test """
                    for key in ret_dict:
                        logging.debug('[+] %s : %s' % (key, ret_dict[key]))
                    """ add file check data """
                    self.add_file_check_data(file_md5, ret_dict["Remote file_id"])
                    return ret_dict["Remote file_id"]
                else:
                    return None
            except Exception as inst:
                exc_msg = traceback.format_exc()
                logging.error(exc_msg)
                logging.error("slave file exists!")
                remote_file_array = remote_file_id.split(".")
                if len(remote_file_array) > 1:
                    remote_file_id = remote_file_array[0] + prefix_name + "." + remote_file_array[1]
                self.add_file_check_data(file_md5, remote_file_id)
                return remote_file_id
        else:
            """ for test """
            logging.debug("file exists!")
            return fileid

    """
    upload_slave_by_buffer:
    """
    '''
        Upload slave file by buffer
        arguments:
        @filebuffer: string
        @remote_file_id: string
        @meta_dict: dictionary e.g.:{
            'ext_name'  : 'jpg',
            'file_size' : '10240B',
            'width'     : '160px',
            'hight'     : '80px'
        }
        @return dictionary {
            'Status'        : 'Upload slave successed.',
            'Local file name' : local_filename,
            'Uploaded size'   : upload_size,
            'Remote file id'  : remote_file_id,
            'Storage IP'      : storage_ip
        }
    '''
    def upload_slave_by_buffer(self, file_md5, filebuffer, remote_file_id, prefix_name=None, \
                               meta_dict = None, check_status=True):
        file_md5 = file_md5 + prefix_name
        fileid = None
        fileid = self.check_file_exists(file_md5)
        if check_status == False:
            self.delete_by_filename(file_md5, fileid)
            fileid = None
        if fileid == None:
            try:
                ret_dict = self.client.upload_slave_by_buffer(filebuffer, remote_file_id, \
                                                            meta_dict=meta_dict, \
                                                            file_ext_name=prefix_name)
                if ret_dict:
                    """ for test """
                    for key in ret_dict:
                        logging.debug('[+] %s : %s' % (key, ret_dict[key]))
                    """ add file check data """
                    self.add_file_check_data(file_md5, ret_dict["Remote file_id"])
                    return ret_dict["Remote file_id"]
                else:
                    return None
            except Exception as inst:

                """ for test """
                logging.error("slave file exists:%s"%inst.args)
                remote_file_array = remote_file_id.split(".")
                if len(remote_file_array) > 1:
                    remote_file_id = remote_file_array[0] + prefix_name + "." + remote_file_array[1]
                self.add_file_check_data(file_md5, remote_file_id)
                return remote_file_id + prefix_name
        else:
            """ for test """
            logging.debug("file exists!")
            return fileid
