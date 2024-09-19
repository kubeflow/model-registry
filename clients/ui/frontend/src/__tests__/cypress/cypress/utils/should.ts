const status = {
  warning: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-warning')).to.eq(true);
  },
  info: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-info')).to.eq(true);
  },
  danger: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-danger')).to.eq(true);
  },
  success: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-success')).to.eq(true);
  },
  custom: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-custom')).to.eq(true);
  },
  error: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-error')).to.eq(true);
  },
  indeterminate: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-indeterminate')).to.eq(true);
  },
};

const expandCollapse = {
  expanded: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-expanded')).to.eq(true);
  },

  collapsed: ($subject: JQuery<HTMLElement>) => {
    expect($subject.hasClass('pf-m-expanded')).to.eq(false);
  },
};

const form = {
  invalid: ($subject: JQuery<HTMLElement>) => {
    expect($subject.attr('aria-invalid')).to.eq('true');
  },
};

const sort = {
  sortAscending: ($subject: JQuery<HTMLElement>) => {
    expect($subject.parents('th').attr('aria-sort')).to.eq('ascending');
  },

  sortDescending: ($subject: JQuery<HTMLElement>) => {
    expect($subject.parents('th').attr('aria-sort')).to.eq('descending');
  },
};

export const be = {
  ...status,
  ...expandCollapse,
  ...form,
  ...sort,
};
