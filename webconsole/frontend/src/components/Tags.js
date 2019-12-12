import React, { Component } from 'react';
import cx from 'classnames';
class Tags extends Component {

  render() {
    let { tags, onAdd, onRemove, theme, fill } = this.props;
    return (
      <div
        className={cx("tagsinput", `tag-${theme}`, {
          'tag-fill': fill === true
        })}
        style={{height: '100%'}}>
        { tags && tags.map(tag => (
          <span className="tag" key={tag.id}>
            <span>{tag.text}</span>&nbsp;<a className="tagsinput-remove-link" onClick={() => onRemove(tag.id)}><i className="fa fa-times"></i></a>
          </span>
        ))}

        <div className="tagsinput-add-container">
          <div className="tagsinput-add"><i className="fa fa-plus"></i></div>
          <input
            defaultValue=""
            style={{color: 'rgb(102, 102, 102)',width: 50}}
            onKeyDown={e => {
              if (e.keyCode === 13 && e.target.value) {
                onAdd(e.target.value);
                e.target.value = '';
              }
              return false;
            }} />
        </div>
      </div>
    )
  }
}

Tags.defaultProps = {
  theme: 'azure',
  fill: false
}

export default Tags;