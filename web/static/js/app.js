$(document).ready(function () {	
	ckb_toggle_all = $('#ck_box_toggle_all');
	ckb_toggle_all.click(function () {
            var $this = $(this);
            $('.bulk_edit_ck_box').prop('checked', $this.is(':checked'));
	});
});