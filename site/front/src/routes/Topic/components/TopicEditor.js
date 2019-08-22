import React from 'react';
import PropTypes from 'prop-types';
import { Card, Button } from 'antd';
import MarkdownEditor from 'components/MarkdownEditor';
import styles from './TopicEditor.scss';

const TopicEditor = ({ value, loading, onEditorChange, onClickReply }) => (
  <Card
    className={styles.container}
    title="添加回复"
    actions={[
      <Button
        type="primary"
        icon="forward"
        loading={loading}
        onClick={onClickReply}
      >
        回复
      </Button>,
    ]}
  >
    <MarkdownEditor value={value} onEditorChange={onEditorChange} />
  </Card>
);

TopicEditor.propTypes = {
  value: PropTypes.string.isRequired,
  loading: PropTypes.bool.isRequired,
  onEditorChange: PropTypes.func.isRequired,
  onClickReply: PropTypes.func.isRequired,
};

export default TopicEditor;
