$(document).ready(function () {	
	ckb_toggle_all = $('#ck_box_toggle_all');
	ckb_toggle_all.click(function () {
            var $this = $(this);
            $('.bulk_edit_ck_box').prop('checked', $this.is(':checked'));
	});
	$('#bulk_export_link0').click(function(e){		
		send_bulk_action('export0','确定导出？')
		return false;
	});
	$('#bulk_export_link1').click(function(e){		
		send_bulk_action('export1','确定导出？')
		return false;
	});

});


function send_bulk_action(action, confirm_message) {        
        var selected_items = $('.bulk_edit_ck_box:checked').clone();
        if (selected_items.length < 1) {
            alert('请选择至少一个项目.');
            return;
        }
        if (confirm_message && confirm(confirm_message) == false) {
            return;
        }
        $('#bulk_action').val(action);
        var hide_div = $('#bulk_hide_div');
        hide_div.find('input:checked').remove();
        hide_div.append(selected_items);
        $('#bulk_form').submit();
}